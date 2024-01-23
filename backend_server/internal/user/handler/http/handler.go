package http

import (
	"backend_server/configs"
	"backend_server/internal/user"
	usr "backend_server/models/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	cfg    *configs.Config
	userUC user.UseCase
}

func NewUserHandler(cfg *configs.Config, clientUC user.UseCase) *UserHandler {
	return &UserHandler{userUC: clientUC, cfg: cfg}
}

func (u *UserHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userReq usr.CreateUserReq
		if err := c.ShouldBindJSON(&userReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := u.userUC.CreateUser(c.Request.Context(), userReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, res)
	}
}

func (u *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userReq usr.LoginUserReq

		if err := c.ShouldBindJSON(&userReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := u.userUC.LoginUser(c.Request.Context(), userReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.SetCookie("jwt", res.AccessToken, 3600, "/", "localhost", false, true)

		resp := usr.LoginUserRes{
			ID:       res.ID,
			Username: res.Username,
		}

		c.JSON(http.StatusOK, resp)
	}
}

func (u *UserHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("jwt", "", -1, "", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	}
}
