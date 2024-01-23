package user

import (
	usr "backend_server/models/user"
	"context"
)

type UseCase interface {
	CreateUser(ctx context.Context, req usr.CreateUserReq) (usr.CreateUserRes, error)
	LoginUser(ctx context.Context, req usr.LoginUserReq) (usr.LoginUserRes, error)
}
