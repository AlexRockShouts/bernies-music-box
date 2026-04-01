package ws

import (
	"encoding/json"
	"log/slog"

	"api-gateway/internal/store"
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	User string
	Send chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// Ignore messages for simplicity, only push
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

type Hub struct {
	clients    map[string][]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *store.Task
	taskStore  *store.TaskStore
	logger     *slog.Logger
}

func NewHub(taskStore *store.TaskStore, logger *slog.Logger) *Hub {
	h := &Hub{
		clients:    make(map[string][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *store.Task),
		taskStore:  taskStore,
		logger:     logger,
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client.User, client)
		case client := <-h.unregister:
			h.removeClient(client.User, client)
		case task := <-h.broadcast:
			h.sendToUser(task.Owner, task)
		}
	}
}

func (h *Hub) addClient(user string, client *Client) {
	h.clients[user] = append(h.clients[user], client)
	// Send initial tasks list
	tasks := h.taskStore.ListByOwner(user)
	data, _ := json.Marshal(tasks)
	select {
	case client.Send <- data:
	default:
		h.logger.Warn("failed to send initial tasks", "user", user)
	}
	h.logger.Info("client registered", "user", user)
}

func (h *Hub) removeClient(user string, client *Client) {
	clients, ok := h.clients[user]
	if !ok {
		return
	}
	for i, c := range clients {
		if c == client {
			h.clients[user] = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	if len(h.clients[user]) == 0 {
		delete(h.clients, user)
	}
	h.logger.Info("client unregistered", "user", user)
}

func (h *Hub) sendToUser(user string, task *store.Task) {
	clients, ok := h.clients[user]
	if !ok {
		return
	}
	data, err := json.Marshal(task)
	if err != nil {
		h.logger.Error("failed to marshal task for WS", "error", err)
		return
	}
	for _, client := range clients {
		select {
		case client.Send <- data:
		default:
			h.logger.Warn("WS send channel full", "user", user)
		}
	}
}

func (h *Hub) Broadcast(task *store.Task) {
	select {
	case h.broadcast <- task:
	default:
		h.logger.Warn("broadcast channel full")
	}
}
