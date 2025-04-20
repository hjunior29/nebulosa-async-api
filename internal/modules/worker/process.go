package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"

	"github.com/hjunior29/nebulosa-async-api/internal/config"
	"github.com/hjunior29/nebulosa-async-api/internal/config/database"
	"github.com/hjunior29/nebulosa-async-api/internal/domain"
)

func buildRequest(task domain.Task) (*http.Request, error) {
	body := bytes.NewReader([]byte{})
	if task.Payload != nil {
		body = bytes.NewReader(task.Payload)
	}

	req, err := http.NewRequest(task.Method, task.Endpoint, body)
	if err != nil {
		return nil, err
	}

	if task.Headers != nil {
		var headers map[string]string
		if err := json.Unmarshal(task.Headers, &headers); err != nil {
			return nil, err
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	return req, nil
}

func ExecuteRequest(task domain.Task) (*http.Response, error) {
	req, err := buildRequest(task)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func StartWorker() {
	go startListener()
	go startPolling()
	go pingAPI()
}

func startListener() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, config.DATABASE_URL)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, "LISTEN new_task")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Worker listening for new tasks...")

	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			log.Printf("notification error: %v", err)
			continue
		}

		go processTaskByID(notification.Payload)
	}
}

func startPolling() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("Worker polling started...")

	for {
		<-ticker.C

		var tasks []domain.Task
		repo := database.NewRepository(&tasks, nil)

		err := repo.FindAllWhere(map[string]interface{}{
			"status":              domain.StatusPending,
			"scheduled_at_time <": time.Now(),
		})
		if err != nil {
			log.Println("Polling error fetching tasks:", err)
			continue
		}

		for _, task := range tasks {
			t := task
			go processTask(t)
		}
	}
}

func processTaskByID(id string) {
	var task domain.Task
	repo := database.NewRepository(&task, nil)
	if err := repo.FindAllWhere(map[string]interface{}{"id": id}); err != nil {
		log.Println("Failed to find task:", err)
		return
	}
	processTask(task)
}

func processTask(task domain.Task) {
	repo := database.NewRepository(&task, nil)

	if task.MaxRetries <= task.Attempts {
		updateTaskStatus(repo, &task, domain.StatusFailed)
		return
	}

	response, err := ExecuteRequest(task)
	if err != nil {
		handleTaskFailure(repo, &task, err, response)
		return
	}

	if response != nil && response.Body != nil {
		response.Body.Close()
	}

	task.StatusCode = response.StatusCode

	updateTaskStatus(repo, &task, domain.StatusSuccess)
}

func updateTaskStatus(repo database.Repository, task *domain.Task, status domain.TaskStatus) {
	task.Status = status
	task.Attempts++
	if err := repo.Save(); err != nil {
		log.Printf("Failed to update task status to %s: %v", status, err)
	}
}

func handleTaskFailure(repo database.Repository, task *domain.Task, err error, response *http.Response) {
	task.LastError = err.Error()
	task.Attempts++
	if response != nil {
		task.StatusCode = response.StatusCode
	}
	if err := repo.Save(); err != nil {
		log.Println("Failed to update task failure:", err)
	}
}

func pingAPI() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	stopAt := time.After(10 * time.Minute)

	for {
		select {
		case <-stopAt:
			log.Println("Stopping API ping after 10 minutes.")
			return
		case <-ticker.C:
			resp, err := http.Get(config.API_URL + "/api/ping")
			if err != nil {
				log.Println("Failed to ping API:", err)
				continue
			}
			if resp.Body != nil {
				resp.Body.Close()
			}
		}
	}
}
