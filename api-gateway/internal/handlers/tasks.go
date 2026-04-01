package handlers

import (
	"net/http"
	"time"

	"api-gateway/internal/services"
	"api-gateway/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
)

func CreateTaskHandler(taskStore *store.TaskStore, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		username := user.(string)
		var req struct {
			Prompt string `json:"prompt" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task := &store.Task{
			ID:        uuid.New().String(),
			Prompt:    req.Prompt,
			Status:    "pending",
			Owner:     username,
			CreatedAt: time.Now(),
		}
		taskStore.Save(task)
		go services.OrchestrateTask(taskStore, task, logger)
		logger.Info("task created", "id", task.ID, "user", username)
		c.JSON(http.StatusAccepted, task)
	}
}

func GetTaskHandler(taskStore *store.TaskStore, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		username := user.(string)
		id := c.Param("id")
		task, ok := taskStore.Get(id)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		if task.Owner != "" && task.Owner != username {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to task"})
			return
		}
		logger.Info("task retrieved", "id", id, "user", username)
		c.JSON(http.StatusOK, task)
	}
}

func ListTasksHandler(taskStore *store.TaskStore, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		username := user.(string)
		list := taskStore.ListByOwner(username)
		logger.Info("tasks listed", "user", username, "count", len(list))
		c.JSON(http.StatusOK, list)
	}
}
