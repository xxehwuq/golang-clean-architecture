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
	repository      repository.UserRepository
	tokensManager   tokens.Manager
	passwordHasher  password.Hasher
	redisDB         *redis.Redis
	refreshTokenTTL time.Duration
}

func NewUserUsecase(repository repository.UserRepository, tokensManager tokens.Manager, passwordHasher password.Hasher, redis *redis.Redis, refreshTokenTTL time.Duration) UserUsecase {
	return &userUsecase{
		repository:      repository,
		tokensManager:   tokensManager,
		passwordHasher:  passwordHasher,
		redisDB:         redis,
		refreshTokenTTL: refreshTokenTTL,
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

	return uc.createSession(ctx, id, nil) // TODO: user permissions
}

func (uc *userUsecase) SignIn(ctx context.Context, input UserSignInInput) (UserTokens, error) {
	user, err := uc.repository.GetByEmail(ctx, input.Email)
	if err != nil {
		return UserTokens{}, err
	}

	if !uc.passwordHasher.Compare(input.Password, user.Password) {
		return UserTokens{}, errors.New("incorrect credentials")
	}

	return uc.createSession(ctx, user.ID, nil) // TODO: user permissions
}

func (uc *userUsecase) RefreshTokens(ctx context.Context, refreshToken, userID string, userPermissions []interface{}) (UserTokens, error) {
	var id string

	if err := uc.redisDB.Get(ctx, refreshToken).Scan(&id); err != nil {
		return UserTokens{}, errors.New("invalid refresh token")
	}

	if userID != id {
		return UserTokens{}, errors.New("invalid refresh token")
	}

	return uc.createSession(ctx, userID, userPermissions)
}

func (uc *userUsecase) createSession(ctx context.Context, userID string, userPermissions []interface{}) (UserTokens, error) {
	accessToken, err := uc.tokensManager.GenerateAccessToken(tokens.UserClaims{
		ID:          userID,
		Permissions: userPermissions,
	})
	if err != nil {
		return UserTokens{}, err
	}

	refreshToken, err := uc.tokensManager.GenerateRefreshToken()
	if err != nil {
		return UserTokens{}, err
	}

	uc.redisDB.Set(ctx, refreshToken, userID, uc.refreshTokenTTL)

	return UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
