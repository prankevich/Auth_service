package usecase

import (
	"auth_service/internal/domain"
	"context"
)

type UserCreater interface {
	CreateUser(ctx context.Context, user domain.User) (err error)
}
