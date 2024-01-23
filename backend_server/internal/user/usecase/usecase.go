package usecase

import (
	"backend_server/configs"
	"backend_server/internal/user"
	usr "backend_server/models/user"
	"backend_server/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

type UserUC struct {
	cfg     *configs.Config
	repo    user.Repository
	timeout time.Duration
}

func NewUserUC(cfg *configs.Config, pgDB user.Repository) *UserUC {
	return &UserUC{
		cfg:     cfg,
		repo:    pgDB,
		timeout: time.Duration(2) * time.Second,
	}
}

func (u *UserUC) CreateUser(ctx context.Context, req usr.CreateUserReq) (usr.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	userData, err := u.repo.GetUserByEmail(ctx, req.Email)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return usr.CreateUserRes{}, err
	}

	if userData.Email != "" {
		return usr.CreateUserRes{}, errors.New("user already exist")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return usr.CreateUserRes{}, err
	}

	req.Password = hashedPassword

	resp, err := u.repo.CreateUserData(ctx, req)
	if err != nil {
		return usr.CreateUserRes{}, err
	}

	return resp.ToServer(), nil
}

func (u *UserUC) LoginUser(ctx context.Context, req usr.LoginUserReq) (usr.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	userRes, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return usr.LoginUserRes{}, err
	}

	err = utils.CheckPassword(req.Password, userRes.Password)
	if err != nil {
		return usr.LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, usr.JWTClaims{
		ID:       strconv.Itoa(int(userRes.ID)),
		Username: userRes.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(userRes.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(u.cfg.SecretKey))
	if err != nil {
		return usr.LoginUserRes{}, err
	}
	return usr.LoginUserRes{
		AccessToken: ss,
		ID:          strconv.Itoa(int(userRes.ID)),
		Username:    userRes.Username,
	}, nil
}
