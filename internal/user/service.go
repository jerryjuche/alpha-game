package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserService struct {
	DBConn *sqlx.DB
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{DBConn: db}

}

func (a *UserService) GetProfile(ctx context.Context) (string ,error) {

}