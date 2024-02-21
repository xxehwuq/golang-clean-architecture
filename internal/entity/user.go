package entity

import (
	"context"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserTokens struct {
	AccessToken  string
	RefreshToken string
}

type UserUsecase interface {
	SignUp(ctx context.Context, u *User) (*UserTokens, error)
	SignIn(ctx context.Context, u *User) (*UserTokens, error)
}
