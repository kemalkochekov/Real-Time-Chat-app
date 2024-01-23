package user

import (
	"github.com/golang-jwt/jwt/v4"
	"strconv"
)

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserReqData struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type CreateUserRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	AccessToken string `json:"-"`
	ID          string `json:"id"`
	Username    string `json:"username"`
}

type JWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (r *CreateUserReq) ToStorage() CreateUserReqData {
	return CreateUserReqData{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}
}

func (r *CreateUserReqData) ToServer() CreateUserRes {
	return CreateUserRes{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}
}
