package handlers

import (
	"net/http"

	"api-gateway/internal/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all for dev
	},
}

func WSHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		username := user.(string)
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		client := &ws.Client{
			Hub:  hub,
			Conn: conn,
			User: username,
			Send: make(chan []byte, 256),
		}
		client.Hub.register <- client
		go client.readPump()
		go client.writePump()
	}
}
