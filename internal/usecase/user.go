package usecase

import (
	"context"
	"errors"
	"github.com/xxehwuq/go-clean-architecture/internal/entity"
	"github.com/xxehwuq/go-clean-architecture/pkg/password"
	"github.com/xxehwuq/go-clean-architecture/pkg/random"
	"github.com/xxehwuq/go-clean-architecture/pkg/tokens"
)

type userUsecase struct {
	repository     entity.UserRepository
	tokensManager  tokens.Manager
	passwordHasher password.Hasher
}

func NewUserUsecase(repository entity.UserRepository, tokensManager tokens.Manager, passwordHasher password.Hasher) entity.UserUsecase {
	return &userUsecase{
		repository:     repository,
		tokensManager:  tokensManager,
		passwordHasher: passwordHasher,
	}
}

func (uc *userUsecase) SignUp(ctx context.Context, u *entity.User) (*entity.UserTokens, error) {
	id, err := random.ID()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := uc.passwordHasher.Hash(u.Password)
	if err != nil {
		return nil, err
	}

	u.ID = id
	u.Password = hashedPassword

	if err = uc.repository.Create(ctx, u); err != nil {
		return nil, err
	}

	accessToken, err := uc.tokensManager.GenerateAccessToken(tokens.UserClaims{
		ID:          u.ID,
		Permissions: nil,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.tokensManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &entity.UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *userUsecase) SignIn(ctx context.Context, u *entity.User) (*entity.UserTokens, error) {
	user, err := uc.repository.GetByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	if !uc.passwordHasher.Compare(u.Password, user.Password) {
		return nil, errors.New("incorrect credentials")
	}

	accessToken, err := uc.tokensManager.GenerateAccessToken(tokens.UserClaims{
		ID:          user.ID,
		Permissions: nil,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.tokensManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &entity.UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
