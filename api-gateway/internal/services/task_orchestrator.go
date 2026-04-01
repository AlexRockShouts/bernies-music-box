package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"api-gateway/internal/store"
)

func OrchestrateTask(taskStore *store.TaskStore, task *store.Task, logger *slog.Logger) {
	task.Status = "processing"
	taskStore.Save(task)

	runpodURL := os.Getenv("RUNPOD_URL")
	if runpodURL == "" {
		runpodURL = "http://localhost:8000/generate"
	}

	payload := map[string]any{
		"lyrics":   task.Prompt,
		"style":    "indie rock, electric guitar, driving drums, male vocal, E minor",
		"duration": 30,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.Error("failed to marshal payload", slog.String("error", err.Error()))
		task.Status = "failed"
		taskStore.Save(task)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", runpodURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("failed to create request", slog.String("error", err.Error()))
		task.Status = "failed"
		taskStore.Save(task)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("AUTH_USERNAME"), os.Getenv("AUTH_PASSWORD"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("runpod request failed", slog.String("error", err.Error()))
		task.Status = "failed"
		taskStore.Save(task)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("runpod non-ok status", slog.Int("status", resp.StatusCode))
		task.Status = "failed"
		taskStore.Save(task)
		return
	}

	var result struct {
		TaskID string `json:"task_id"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("failed to decode response", slog.String("error", err.Error()))
		task.Status = "failed"
		taskStore.Save(task)
		return
	}

	go PollBackend(taskStore, task, result.TaskID, logger)
}

func PollBackend(taskStore *store.TaskStore, task *store.Task, backendTaskID string, logger *slog.Logger) {
	runpodURL := os.Getenv("RUNPOD_URL")
	if runpodURL == "" {
		runpodURL = "http://localhost:8000"
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.NewTimer(10 * time.Minute)
	defer timeout.Stop()

	for {
		select {
		case <-ticker.C:
			url := fmt.Sprintf("%s/task/%s", runpodURL, backendTaskID)
			req, _ := http.NewRequest("GET", url, nil)
			req.SetBasicAuth(os.Getenv("AUTH_USERNAME"), os.Getenv("AUTH_PASSWORD"))
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				logger.Warn("poll backend failed", slog.String("error", err.Error()))
				continue
			}
			defer resp.Body.Close()

			var res struct {
				Status   string `json:"status"`
				AudioURL string `json:"audio_url"`
				MIDIURL  string `json:"midi_url"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
				logger.Warn("poll decode failed", slog.String("error", err.Error()))
				continue
			}

			if res.Status == "completed" {
				task.Status = "completed"
				task.ResultURL = res.AudioURL
				taskStore.Save(task)
				return
			} else if res.Status == "failed" {
				task.Status = "failed"
				taskStore.Save(task)
				return
			}
		case <-timeout.C:
			task.Status = "failed"
			taskStore.Save(task)
			return
		}
	}
}
