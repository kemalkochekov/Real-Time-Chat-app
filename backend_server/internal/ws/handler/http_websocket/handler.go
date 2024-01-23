package http_websocket

import (
	ws2 "backend_server/internal/ws"
	"backend_server/models/web_socket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type Handler struct {
	hub *ws2.Hub
}

func NewHandler(h *ws2.Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

func (h *Handler) CreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req web_socket.CreateRoomReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.hub.Rooms[req.RoomID] = &ws2.Room{
			RoomID:  req.RoomID,
			Name:    req.Name,
			Clients: make(map[string]*ws2.Client),
		}
		c.JSON(http.StatusOK, req)
	}
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		roomID := c.Param("roomId")
		clientID := c.Query("userId")
		username := c.Query("username")

		cl := &ws2.Client{
			Conn:     conn,
			Message:  make(chan *web_socket.Message, 10),
			ID:       clientID,
			RoomID:   roomID,
			Username: username,
		}

		m := &web_socket.Message{
			Content:  "A new user has joined the room",
			RoomID:   roomID,
			Username: username,
		}

		// Register a new client through the register channel
		h.hub.Register <- cl
		// Broadcast that message
		h.hub.Broadcast <- m

		go cl.WriteMessage()
		cl.ReadMessage(h.hub)
	}
}

func (h *Handler) GetRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms := make([]web_socket.RoomRes, 0)

		for _, r := range h.hub.Rooms {
			rooms = append(rooms, web_socket.RoomRes{
				ID:   r.RoomID,
				Name: r.Name,
			})
		}

		c.JSON(http.StatusOK, rooms)
	}
}

func (h *Handler) GetClients() gin.HandlerFunc {
	return func(c *gin.Context) {
		var clients []web_socket.ClientRes
		roomID := c.Param("roomId")

		if _, ok := h.hub.Rooms[roomID]; !ok {
			clients = make([]web_socket.ClientRes, 0)
			c.JSON(http.StatusOK, clients)
			return
		}

		for _, c := range h.hub.Rooms[roomID].Clients {
			clients = append(clients, web_socket.ClientRes{
				ID:       c.ID,
				Username: c.Username,
			})
		}

		c.JSON(http.StatusOK, clients)
	}
}
