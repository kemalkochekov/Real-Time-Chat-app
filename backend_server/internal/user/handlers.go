package user

import "github.com/gin-gonic/gin"

type Handlers interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	Logout() gin.HandlerFunc
}
