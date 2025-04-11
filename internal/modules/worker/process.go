package worker

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	timer := 1 * time.Second
	taskTicker := time.NewTicker(timer)
	defer taskTicker.Stop()

	pingTicker := time.NewTicker(10 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-taskTicker.C:
			var pendingTasks []domain.Task
			repo := database.NewRepository(&pendingTasks, nil)

			if err := repo.FindAllWhere(map[string]interface{}{"status": domain.StatusPending}); err != nil {
				log.Println("Error fetching pending tasks:", err)
				continue
			}

			if len(pendingTasks) == 0 {
				taskTicker.Stop()
				timer = time.Duration(float64(timer) * 1.2)

				if timer > 10*time.Second {
					timer = 1 * time.Second
				}
				taskTicker = time.NewTicker(timer)
				log.Println("No pending tasks, sleeping for", timer)
				continue
			}

			timer = 1 * time.Second
			taskTicker.Reset(timer)

			for _, task := range pendingTasks {
				taskRepo := database.NewRepository(&task, nil)

				if task.MaxRetries <= task.Attempts {
					updateTaskStatus(taskRepo, &task, domain.StatusFailed)
					continue
				}

				response, err := ExecuteRequest(task)
				if err != nil {
					handleTaskFailure(taskRepo, &task, err, response)
					continue
				}

				if response != nil && response.Body != nil {
					response.Body.Close()
				}

				updateTaskStatus(taskRepo, &task, domain.StatusSuccess)
			}
		case <-pingTicker.C:
			_, err := http.Get(config.API_URL + "/api/ping")
			if err != nil {
				log.Println("Error pinging worker:", err)
				continue
			}
		}
	}
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
