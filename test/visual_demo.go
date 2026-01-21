package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL  = "http://localhost:8080/api"
	username = "admin"
	password = "adajhsvdgahsvdaghsgds"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

type TaskRequest struct {
	Endpoint    string            `json:"endpoint"`
	Method      string            `json:"method"`
	Payload     map[string]any    `json:"payload"`
	Headers     map[string]string `json:"headers"`
	Type        string            `json:"type"`
	MaxRetries  int               `json:"maxRetries"`
	ScheduledAt string            `json:"scheduledAt,omitempty"`
}

type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

var (
	httpClient = &http.Client{Timeout: 30 * time.Second}
	token      string
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║     NEBULOSA ASYNC API - VISUAL TEST SUITE                   ║")
	fmt.Println("║     Este teste cria várias tasks para visualização           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Login
	fmt.Println("🔐 Fazendo login...")
	if err := login(); err != nil {
		fmt.Printf("❌ Erro no login: %v\n", err)
		return
	}
	fmt.Println("✅ Login realizado com sucesso!")
	fmt.Println()

	// ========== PRIMEIRO: TODAS AS TASKS EM PARALELO (BURST) ==========
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🚀 FASE 1: BURST INICIAL - Criando 30 tasks SIMULTANEAMENTE")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	createAllTasksInParallel()

	time.Sleep(3 * time.Second)

	// ========== DEPOIS: TASKS EM SEQUÊNCIA ==========
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📦 FASE 2: Tasks de SUCESSO em sequência (uma por uma)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	createSuccessTasks()

	time.Sleep(2 * time.Second)

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💥 FASE 3: Tasks de FALHA em sequência (uma por uma)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	createFailureTasks()

	time.Sleep(2 * time.Second)

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("⏰ FASE 4: Tasks AGENDADAS em sequência (uma por uma)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	createScheduledTasks()

	time.Sleep(2 * time.Second)

	// Loop contínuo criando tasks aleatórias por 2 minutos
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔄 FASE 5: Loop contínuo de 2 minutos criando tasks aleatórias")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	continuousTaskCreation(2 * time.Minute)

	// Listar todas as tasks no final
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 FASE FINAL: Listando todas as tasks")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	listAllTasks()

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              TESTE FINALIZADO COM SUCESSO!                   ║")
	fmt.Println("║     Verifique o frontend para visualizar as tasks           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

func login() error {
	loginReq := LoginRequest{
		Username: username,
		Password: password,
	}

	body, _ := json.Marshal(loginReq)
	resp, err := httpClient.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return err
	}

	if loginResp.Status != 200 {
		return fmt.Errorf("login failed: %s", loginResp.Message)
	}

	token = loginResp.Data.Token
	return nil
}

func createTask(task TaskRequest) error {
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", baseURL+"/task", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("   → Task criada: %s (Status: %d)\n", task.Type, resp.StatusCode)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to create task: %s", string(respBody))
	}

	return nil
}

func createSuccessTasks() {
	tasks := []TaskRequest{
		{
			Endpoint:   "https://httpbin.org/post",
			Method:     "POST",
			Payload:    map[string]any{"message": "Hello from Nebulosa!", "timestamp": time.Now().Unix()},
			Headers:    map[string]string{"X-Custom-Header": "nebulosa-test"},
			Type:       "webhook-notification",
			MaxRetries: 3,
		},
		{
			Endpoint:   "https://httpbin.org/get",
			Method:     "GET",
			Payload:    map[string]any{},
			Headers:    map[string]string{"Accept": "application/json"},
			Type:       "health-check",
			MaxRetries: 1,
		},
		{
			Endpoint:   "https://jsonplaceholder.typicode.com/posts",
			Method:     "POST",
			Payload:    map[string]any{"title": "Nebulosa Test", "body": "Testing async queue", "userId": 1},
			Headers:    map[string]string{"Content-Type": "application/json"},
			Type:       "api-integration",
			MaxRetries: 2,
		},
		{
			Endpoint:   "https://httpbin.org/put",
			Method:     "PUT",
			Payload:    map[string]any{"updated": true, "version": "2.0"},
			Headers:    map[string]string{},
			Type:       "data-sync",
			MaxRetries: 3,
		},
		{
			Endpoint:   "https://httpbin.org/patch",
			Method:     "PATCH",
			Payload:    map[string]any{"partial": "update"},
			Headers:    map[string]string{},
			Type:       "partial-update",
			MaxRetries: 2,
		},
	}

	for _, task := range tasks {
		if err := createTask(task); err != nil {
			fmt.Printf("   ❌ Erro: %v\n", err)
		}
		time.Sleep(500 * time.Millisecond) // Delay para visualização
	}
}

func createFailureTasks() {
	tasks := []TaskRequest{
		{
			Endpoint:   "https://invalid-domain-that-does-not-exist.com/api",
			Method:     "POST",
			Payload:    map[string]any{"test": "should fail"},
			Headers:    map[string]string{},
			Type:       "invalid-domain",
			MaxRetries: 2,
		},
		{
			Endpoint:   "https://httpbin.org/status/500",
			Method:     "GET",
			Payload:    map[string]any{},
			Headers:    map[string]string{},
			Type:       "server-error-500",
			MaxRetries: 3,
		},
		{
			Endpoint:   "https://httpbin.org/status/404",
			Method:     "GET",
			Payload:    map[string]any{},
			Headers:    map[string]string{},
			Type:       "not-found-404",
			MaxRetries: 2,
		},
		{
			Endpoint:   "https://httpbin.org/status/401",
			Method:     "POST",
			Payload:    map[string]any{"auth": "missing"},
			Headers:    map[string]string{},
			Type:       "unauthorized-401",
			MaxRetries: 1,
		},
		{
			Endpoint:   "https://httpbin.org/delay/30",
			Method:     "GET",
			Payload:    map[string]any{},
			Headers:    map[string]string{},
			Type:       "timeout-expected",
			MaxRetries: 1,
		},
	}

	for _, task := range tasks {
		if err := createTask(task); err != nil {
			fmt.Printf("   ❌ Erro: %v\n", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func createScheduledTasks() {
	tasks := []TaskRequest{
		{
			Endpoint:    "https://httpbin.org/post",
			Method:      "POST",
			Payload:     map[string]any{"scheduled": true, "delay": "30s"},
			Headers:     map[string]string{},
			Type:        "scheduled-30s",
			MaxRetries:  2,
			ScheduledAt: "30s",
		},
		{
			Endpoint:    "https://httpbin.org/post",
			Method:      "POST",
			Payload:     map[string]any{"scheduled": true, "delay": "1m"},
			Headers:     map[string]string{},
			Type:        "scheduled-1min",
			MaxRetries:  2,
			ScheduledAt: "1m",
		},
		{
			Endpoint:    "https://httpbin.org/post",
			Method:      "POST",
			Payload:     map[string]any{"scheduled": true, "delay": "2m"},
			Headers:     map[string]string{},
			Type:        "scheduled-2min",
			MaxRetries:  2,
			ScheduledAt: "2m",
		},
		{
			Endpoint:    "https://httpbin.org/post",
			Method:      "POST",
			Payload:     map[string]any{"scheduled": true, "delay": "5m"},
			Headers:     map[string]string{},
			Type:        "scheduled-5min",
			MaxRetries:  2,
			ScheduledAt: "5m",
		},
	}

	for _, task := range tasks {
		if err := createTask(task); err != nil {
			fmt.Printf("   ❌ Erro: %v\n", err)
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func createAllTasksInParallel() {
	var wg sync.WaitGroup

	// Todas as tasks que serão criadas simultaneamente
	allTasks := []TaskRequest{
		// Tasks de sucesso
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"message": "Hello!", "type": "success"}, Headers: map[string]string{"X-Type": "success"}, Type: "burst-webhook-1", MaxRetries: 3},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"message": "Hello!", "type": "success"}, Headers: map[string]string{"X-Type": "success"}, Type: "burst-webhook-2", MaxRetries: 3},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"message": "Hello!", "type": "success"}, Headers: map[string]string{"X-Type": "success"}, Type: "burst-webhook-3", MaxRetries: 3},
		{Endpoint: "https://httpbin.org/get", Method: "GET", Payload: map[string]any{}, Headers: map[string]string{}, Type: "burst-health-check-1", MaxRetries: 1},
		{Endpoint: "https://httpbin.org/get", Method: "GET", Payload: map[string]any{}, Headers: map[string]string{}, Type: "burst-health-check-2", MaxRetries: 1},
		{Endpoint: "https://jsonplaceholder.typicode.com/posts", Method: "POST", Payload: map[string]any{"title": "Burst", "body": "Test"}, Headers: map[string]string{}, Type: "burst-api-call-1", MaxRetries: 2},
		{Endpoint: "https://jsonplaceholder.typicode.com/posts", Method: "POST", Payload: map[string]any{"title": "Burst", "body": "Test"}, Headers: map[string]string{}, Type: "burst-api-call-2", MaxRetries: 2},
		{Endpoint: "https://httpbin.org/put", Method: "PUT", Payload: map[string]any{"updated": true}, Headers: map[string]string{}, Type: "burst-data-sync-1", MaxRetries: 2},
		{Endpoint: "https://httpbin.org/put", Method: "PUT", Payload: map[string]any{"updated": true}, Headers: map[string]string{}, Type: "burst-data-sync-2", MaxRetries: 2},
		{Endpoint: "https://httpbin.org/patch", Method: "PATCH", Payload: map[string]any{"partial": true}, Headers: map[string]string{}, Type: "burst-partial-update", MaxRetries: 2},

		// Tasks que vão falhar
		{Endpoint: "https://invalid-domain-xyz.com/api", Method: "POST", Payload: map[string]any{"fail": true}, Headers: map[string]string{}, Type: "burst-invalid-domain-1", MaxRetries: 2},
		{Endpoint: "https://invalid-domain-abc.com/api", Method: "POST", Payload: map[string]any{"fail": true}, Headers: map[string]string{}, Type: "burst-invalid-domain-2", MaxRetries: 2},
		{Endpoint: "https://httpbin.org/status/500", Method: "GET", Payload: map[string]any{}, Headers: map[string]string{}, Type: "burst-error-500-1", MaxRetries: 3},
		{Endpoint: "https://httpbin.org/status/500", Method: "GET", Payload: map[string]any{}, Headers: map[string]string{}, Type: "burst-error-500-2", MaxRetries: 3},
		{Endpoint: "https://httpbin.org/status/404", Method: "GET", Payload: map[string]any{}, Headers: map[string]string{}, Type: "burst-error-404", MaxRetries: 2},
		{Endpoint: "https://httpbin.org/status/401", Method: "POST", Payload: map[string]any{}, Headers: map[string]string{}, Type: "burst-error-401", MaxRetries: 1},
		{Endpoint: "https://httpbin.org/status/503", Method: "GET", Payload: map[string]any{}, Headers: map[string]string{}, Type: "burst-error-503", MaxRetries: 2},

		// Tasks agendadas
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"scheduled": true}, Headers: map[string]string{}, Type: "burst-scheduled-15s", MaxRetries: 2, ScheduledAt: "15s"},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"scheduled": true}, Headers: map[string]string{}, Type: "burst-scheduled-30s", MaxRetries: 2, ScheduledAt: "30s"},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"scheduled": true}, Headers: map[string]string{}, Type: "burst-scheduled-45s", MaxRetries: 2, ScheduledAt: "45s"},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"scheduled": true}, Headers: map[string]string{}, Type: "burst-scheduled-1m", MaxRetries: 2, ScheduledAt: "1m"},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"scheduled": true}, Headers: map[string]string{}, Type: "burst-scheduled-90s", MaxRetries: 2, ScheduledAt: "90s"},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"scheduled": true}, Headers: map[string]string{}, Type: "burst-scheduled-2m", MaxRetries: 2, ScheduledAt: "2m"},

		// Mais tasks variadas
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"type": "email"}, Headers: map[string]string{}, Type: "burst-email-notification", MaxRetries: 3},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"type": "sms"}, Headers: map[string]string{}, Type: "burst-sms-alert", MaxRetries: 3},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"type": "push"}, Headers: map[string]string{}, Type: "burst-push-notification", MaxRetries: 3},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"type": "report"}, Headers: map[string]string{}, Type: "burst-report-generation", MaxRetries: 2},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"type": "backup"}, Headers: map[string]string{}, Type: "burst-backup-task", MaxRetries: 2},
		{Endpoint: "https://httpbin.org/post", Method: "POST", Payload: map[string]any{"type": "cleanup"}, Headers: map[string]string{}, Type: "burst-cleanup-job", MaxRetries: 1},
		{Endpoint: "https://httpbin.org/delay/5", Method: "GET", Payload: map[string]any{}, Headers: map[string]string{}, Type: "burst-slow-request", MaxRetries: 1},
	}

	fmt.Printf("   Disparando %d tasks SIMULTANEAMENTE...\n", len(allTasks))

	for _, task := range allTasks {
		wg.Add(1)
		go func(t TaskRequest) {
			defer wg.Done()
			createTask(t)
		}(task)
	}

	wg.Wait()
	fmt.Printf("   ✅ Todas as %d tasks foram criadas de uma vez!\n", len(allTasks))
}

func continuousTaskCreation(duration time.Duration) {
	endTime := time.Now().Add(duration)
	taskCounter := 0

	taskTypes := []string{
		"email-notification",
		"sms-alert",
		"push-notification",
		"data-export",
		"report-generation",
		"file-processing",
		"image-resize",
		"video-transcode",
		"pdf-generation",
		"backup-task",
		"cleanup-job",
		"analytics-event",
		"audit-log",
		"cache-invalidation",
		"index-rebuild",
	}

	endpoints := []string{
		"https://httpbin.org/post",
		"https://httpbin.org/put",
		"https://jsonplaceholder.typicode.com/posts",
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/500", // Alguns vão falhar
		"https://httpbin.org/delay/5",
	}

	methods := []string{"POST", "PUT", "PATCH"}

	scheduledOptions := []string{"", "", "", "10s", "20s", "30s", "1m"} // Maioria imediata

	fmt.Printf("   Iniciando loop de %v...\n", duration)
	fmt.Println("   [Pressione Ctrl+C para interromper]")
	fmt.Println()

	for time.Now().Before(endTime) {
		taskCounter++

		taskType := taskTypes[rand.Intn(len(taskTypes))]
		endpoint := endpoints[rand.Intn(len(endpoints))]
		method := methods[rand.Intn(len(methods))]
		scheduled := scheduledOptions[rand.Intn(len(scheduledOptions))]

		task := TaskRequest{
			Endpoint: endpoint,
			Method:   method,
			Payload: map[string]any{
				"taskNumber": taskCounter,
				"taskType":   taskType,
				"timestamp":  time.Now().Format(time.RFC3339),
				"randomData": rand.Int63(),
			},
			Headers:     map[string]string{"X-Task-Number": fmt.Sprintf("%d", taskCounter)},
			Type:        fmt.Sprintf("%s-%d", taskType, taskCounter),
			MaxRetries:  rand.Intn(3) + 1,
			ScheduledAt: scheduled,
		}

		if err := createTask(task); err != nil {
			fmt.Printf("   ❌ Task #%d falhou: %v\n", taskCounter, err)
		}

		// Intervalo aleatório entre 2-5 segundos
		sleepTime := time.Duration(2000+rand.Intn(3000)) * time.Millisecond
		time.Sleep(sleepTime)

		// Mostra progresso a cada 10 tasks
		if taskCounter%10 == 0 {
			remaining := time.Until(endTime)
			fmt.Printf("   📊 Progresso: %d tasks criadas | Tempo restante: %v\n", taskCounter, remaining.Round(time.Second))
		}
	}

	fmt.Printf("\n   ✅ Loop finalizado! Total de tasks criadas: %d\n", taskCounter)
}

func listAllTasks() {
	req, _ := http.NewRequest("GET", baseURL+"/task", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("   ❌ Erro ao listar tasks: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("   ❌ Erro ao decodificar resposta: %v\n", err)
		return
	}

	if tasks, ok := apiResp.Data.([]interface{}); ok {
		statusCount := make(map[string]int)

		for _, t := range tasks {
			if task, ok := t.(map[string]interface{}); ok {
				if status, ok := task["status"].(string); ok {
					statusCount[status]++
				}
			}
		}

		fmt.Printf("\n   📊 RESUMO DAS TASKS:\n")
		fmt.Printf("   ────────────────────────────\n")
		fmt.Printf("   Total de tasks: %d\n", len(tasks))
		fmt.Println()
		for status, count := range statusCount {
			emoji := getStatusEmoji(status)
			fmt.Printf("   %s %s: %d\n", emoji, status, count)
		}
		fmt.Printf("   ────────────────────────────\n")
	}
}

func getStatusEmoji(status string) string {
	switch status {
	case "pending":
		return "⏳"
	case "processing":
		return "🔄"
	case "success":
		return "✅"
	case "failed":
		return "❌"
	default:
		return "❓"
	}
}
