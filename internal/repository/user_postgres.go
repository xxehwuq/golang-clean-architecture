package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/xxehwuq/go-clean-architecture/internal/entity"
	"github.com/xxehwuq/go-clean-architecture/pkg/postgres"
)

type userRepository struct {
	*postgres.Postgres
	tableName string
}

func NewUserRepository(postgres *postgres.Postgres, tableName string) entity.UserRepository {
	return &userRepository{
		Postgres:  postgres,
		tableName: tableName,
	}
}

func (r *userRepository) Create(ctx context.Context, u *entity.User) error {
	sql, args, err := r.Builder.
		Insert(r.tableName).
		Columns("id", "name", "email", "password", "created_at", "updated_at").
		Values(u.ID, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	sql, args, err := r.Builder.
		Select("id", "name", "email", "password", "created_at", "updated_at").
		From(r.tableName).
		Where(squirrel.Eq{"email": email}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user entity.User

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, err
}
