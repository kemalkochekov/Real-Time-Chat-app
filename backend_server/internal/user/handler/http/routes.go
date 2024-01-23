package http

import (
	"backend_server/internal/user"
	"github.com/gin-gonic/gin"
)

func MapUserRoutes(group *gin.RouterGroup, h user.Handlers) {
	group.POST("/signup", h.Register())
	group.POST("/login", h.Login())
	group.GET("/logout", h.Logout())
}
