package http_websocket

import (
	"backend_server/internal/ws"
	"github.com/gin-gonic/gin"
)

func MapUserRoutes(group *gin.RouterGroup, h ws.WebSocketHandler) {
	group.POST("/createRoom", h.CreateRoom())
	group.GET("/joinRoom/:roomId", h.JoinRoom())
	group.GET("/getRooms", h.GetRooms())
	group.GET("/getClients/:roomId", h.GetClients())
}
