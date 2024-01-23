package user

import (
	usr "backend_server/models/user"
	"context"
)

type Repository interface {
	CreateUserData(ctx context.Context, user usr.CreateUserReq) (usr.CreateUserReqData, error)
	GetUserByEmail(ctx context.Context, email string) (usr.CreateUserReqData, error)
}
