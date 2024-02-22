package usecase

import "context"

type UserSignUpInput struct {
	Name     string
	Email    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type UserTokens struct {
	AccessToken  string
	RefreshToken string
}

type UserUsecase interface {
	SignUp(ctx context.Context, input UserSignUpInput) (UserTokens, error)
	SignIn(ctx context.Context, input UserSignInInput) (UserTokens, error)
	RefreshTokens(ctx context.Context, refreshToken, userID string) (UserTokens, error)
}
