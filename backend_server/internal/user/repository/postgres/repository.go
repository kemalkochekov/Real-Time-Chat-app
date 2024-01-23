package postgres

import (
	"backend_server/configs"
	usr "backend_server/models/user"
	"backend_server/pkg/connection/postgres"
	"context"
)

type UserRepo struct {
	cfg    *configs.Config
	driver postgres.DBops
}

func NewUserRepo(cfg *configs.Config, drv postgres.DBops) *UserRepo {
	return &UserRepo{
		cfg:    cfg,
		driver: drv,
	}
}

func (u *UserRepo) CreateUserData(ctx context.Context, user usr.CreateUserReq) (usr.CreateUserReqData, error) {
	var lastInsertedID int64

	userData := user.ToStorage()
	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"

	err := u.driver.QueryRowContext(ctx, query, userData.Username, userData.Password, userData.Email).Scan(&lastInsertedID)
	if err != nil {
		return usr.CreateUserReqData{}, err
	}

	userData.ID = lastInsertedID

	return userData, nil
}

func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (usr.CreateUserReqData, error) {
	var req usr.CreateUserReqData

	query := "SELECT id, email, username, password FROM users WHERE email = $1"

	err := u.driver.QueryRowContext(ctx, query, email).Scan(&req.ID, &req.Email, &req.Username, &req.Password)
	if err != nil {
		return usr.CreateUserReqData{}, err
	}

	return req, nil
}
