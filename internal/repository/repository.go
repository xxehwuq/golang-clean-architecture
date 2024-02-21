package repository

import (
	"context"
	"github.com/xxehwuq/go-clean-architecture/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) error
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}
