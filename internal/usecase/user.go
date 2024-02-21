package usecase

import (
	"context"
	"errors"
	"github.com/xxehwuq/go-clean-architecture/internal/entity"
	"github.com/xxehwuq/go-clean-architecture/internal/repository"
	"github.com/xxehwuq/go-clean-architecture/pkg/password"
	"github.com/xxehwuq/go-clean-architecture/pkg/random"
	"github.com/xxehwuq/go-clean-architecture/pkg/redis"
	"github.com/xxehwuq/go-clean-architecture/pkg/tokens"
	"time"
)

type userUsecase struct {
	repository     repository.UserRepository
	tokensManager  tokens.Manager
	passwordHasher password.Hasher
	redisDB        *redis.Redis
}

func NewUserUsecase(repository repository.UserRepository, tokensManager tokens.Manager, passwordHasher password.Hasher, redis *redis.Redis) UserUsecase {
	return &userUsecase{
		repository:     repository,
		tokensManager:  tokensManager,
		passwordHasher: passwordHasher,
		redisDB:        redis,
	}
}

func (uc *userUsecase) SignUp(ctx context.Context, input UserSignUpInput) (UserTokens, error) {
	id, err := random.ID()
	if err != nil {
		return UserTokens{}, err
	}

	hashedPassword, err := uc.passwordHasher.Hash(input.Password)
	if err != nil {
		return UserTokens{}, err
	}

	if err = uc.repository.Create(ctx, entity.User{
		ID:        id,
		Name:      input.Name,
		Email:     input.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		return UserTokens{}, err
	}

	accessToken, err := uc.tokensManager.GenerateAccessToken(tokens.UserClaims{
		ID:          id,
		Permissions: nil,
	})
	if err != nil {
		return UserTokens{}, err
	}

	refreshToken, err := uc.tokensManager.GenerateRefreshToken()
	if err != nil {
		return UserTokens{}, err
	}

	return UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *userUsecase) SignIn(ctx context.Context, input UserSignInInput) (UserTokens, error) {
	user, err := uc.repository.GetByEmail(ctx, input.Email)
	if err != nil {
		return UserTokens{}, err
	}

	if !uc.passwordHasher.Compare(input.Password, user.Password) {
		return UserTokens{}, errors.New("incorrect credentials")
	}

	accessToken, err := uc.tokensManager.GenerateAccessToken(tokens.UserClaims{
		ID:          user.ID,
		Permissions: nil,
	})
	if err != nil {
		return UserTokens{}, err
	}

	refreshToken, err := uc.tokensManager.GenerateRefreshToken()
	if err != nil {
		return UserTokens{}, err
	}

	return UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
