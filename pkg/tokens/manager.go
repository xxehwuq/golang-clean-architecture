package tokens

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type Manager interface {
	GenerateAccessToken(userClaims UserClaims) (string, error)
	GenerateRefreshToken() (string, error)
	ParseUserClaims(accessToken string) (UserClaims, error)
}

type UserClaims struct {
	ID string
}

type manager struct {
	signingKey     string
	accessTokenTTL time.Duration
}

func New(signingKey string, accessTokenTTL time.Duration) Manager {
	return &manager{
		signingKey:     signingKey,
		accessTokenTTL: accessTokenTTL,
	}
}

func (m *manager) GenerateAccessToken(userClaims UserClaims) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(m.accessTokenTTL).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = userClaims.ID

	tokenString, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *manager) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (m *manager) ParseUserClaims(accessToken string) (UserClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return UserClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return UserClaims{}, fmt.Errorf("error getting user claims from access token")
	}

	return UserClaims{
		ID: claims["sub"].(string),
	}, nil
}
