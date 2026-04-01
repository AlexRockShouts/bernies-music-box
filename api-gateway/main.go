package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Task struct {
	ID        string    `json:"id"`
	Prompt    string    `json:"prompt"`
	Status    string    `json:"status"`
	ResultURL string    `json:"result_url,omitzero"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskStore struct {
	sync.RWMutex
	tasks map[string]*Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]*Task),
	}
}

func (s *TaskStore) Save(t *Task) {
	s.Lock()
	defer s.Unlock()
	t.UpdatedAt = time.Now()
	s.tasks[t.ID] = t
}

func (s *TaskStore) Get(id string) (*Task, bool) {
	s.RLock()
	defer s.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *TaskStore) List() []*Task {
	s.RLock()
	defer s.RUnlock()
	list := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	return list
}

var store = NewTaskStore()

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: %v", err)
	}
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/tasks", createTask)
	r.GET("/tasks/:id", getTask)
	r.GET("/history", listTasks)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func createTask(c *gin.Context) {
	var req struct {
		Prompt string `json:"prompt" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &Task{
		ID:        uuid.New().String(),
		Prompt:    req.Prompt,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	store.Save(task)

	// Orchestration: Async call to RunPod
	go orchestrateTask(task)

	c.JSON(http.StatusAccepted, task)
}

func getTask(c *gin.Context) {
	id := c.Param("id")
	task, ok := store.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func listTasks(c *gin.Context) {
	c.JSON(http.StatusOK, store.List())
}

func orchestrateTask(task *Task) {
	task.Status = "processing"
	store.Save(task)

	runpodURL := os.Getenv("RUNPOD_URL")
	if runpodURL == "" {
		runpodURL = "http://localhost:8000/generate" // Local FastAPI fallback
	}

	payload := map[string]any{
		"lyrics":   task.Prompt,
		"style":    "indie rock, electric guitar, driving drums, male vocal, E minor",
		"duration": 30,
	}
	jsonData, _ := json.Marshal(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "POST", runpodURL, bytes.NewBuffer(jsonData))
	req.SetBasicAuth(
		os.Getenv("AUTH_USERNAME"),
		os.Getenv("AUTH_PASSWORD"),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling RunPod: %v", err)
		task.Status = "failed"
		store.Save(task)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("RunPod returned status: %d", resp.StatusCode)
		task.Status = "failed"
		store.Save(task)
		return
	}

	var result struct {
		TaskID string `json:"task_id"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding RunPod response: %v", err)
		task.Status = "failed"
		store.Save(task)
		return
	}

	// Since the backend is async, we need to poll or just update status
	task.Status = "processing"
	// For simplicity in this mock, we'll wait or just update the UI that it's processing
	// In a real app, the frontend would poll /tasks/:id which would poll the backend
	store.Save(task)

	// Polling simulation for this mock
	go pollBackend(task, result.TaskID)
}

func pollBackend(task *Task, backendTaskID string) {
	runpodURL := os.Getenv("RUNPOD_URL")
	if runpodURL == "" {
		runpodURL = "http://localhost:8000"
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			url := fmt.Sprintf("%s/task/%s", runpodURL, backendTaskID)
			req, _ := http.NewRequest("GET", url, nil)
			req.SetBasicAuth(os.Getenv("AUTH_USERNAME"), os.Getenv("AUTH_PASSWORD"))

			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Polling error: %v", err)
				continue
			}

			var res struct {
				Status   string `json:"status"`
				AudioURL string `json:"audio_url"`
				MIDIURL  string `json:"midi_url"`
			}
			json.NewDecoder(resp.Body).Decode(&res)
			resp.Body.Close()

			if res.Status == "completed" {
				task.Status = "completed"
				// Combine audio and midi URLs if needed or just pick audio
				// For the UI we previously had ResultURL
				task.ResultURL = res.AudioURL
				store.Save(task)
				return
			} else if res.Status == "failed" {
				task.Status = "failed"
				store.Save(task)
				return
			}
		case <-time.After(10 * time.Minute):
			task.Status = "failed"
			store.Save(task)
			return
		}
	}
}
