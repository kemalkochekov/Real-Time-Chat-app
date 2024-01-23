package ws

import "github.com/gin-gonic/gin"

type WebSocketHandler interface {
	CreateRoom() gin.HandlerFunc
	JoinRoom() gin.HandlerFunc
	GetRooms() gin.HandlerFunc
	GetClients() gin.HandlerFunc
}
