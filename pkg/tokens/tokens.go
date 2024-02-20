package tokens

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type Interface interface {
	GenerateAccessToken(userClaims UserClaims) (AccessToken, error)
	GenerateRefreshToken() (RefreshToken, error)
	ParseUserClaims(accessToken AccessToken) (UserClaims, error)
	ValidateAccessToken() bool
}

type AccessToken string

type RefreshToken string

type UserClaims struct {
	ID          string
	Permissions []interface{}
}

type tokens struct {
	signingKey      string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func New(signingKey string, accessTokenTTL, refreshTokenTTL time.Duration) Interface {
	return &tokens{
		signingKey:      signingKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (t *tokens) GenerateAccessToken(userClaims UserClaims) (AccessToken, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(t.accessTokenTTL).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = userClaims.ID
	claims["permissions"] = userClaims.Permissions

	tokenString, err := token.SignedString([]byte(t.signingKey))
	if err != nil {
		return "", err
	}

	return AccessToken(tokenString), nil
}

func (t *tokens) GenerateRefreshToken() (RefreshToken, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return RefreshToken(fmt.Sprintf("%x", b)), nil
}

func (t *tokens) ParseUserClaims(accessToken AccessToken) (UserClaims, error) {
	token, err := jwt.Parse(string(accessToken), func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(t.signingKey), nil
	})
	if err != nil {
		return UserClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return UserClaims{}, fmt.Errorf("error getting user claims from access token")
	}

	return UserClaims{
		ID:          claims["sub"].(string),
		Permissions: claims["permissions"].([]interface{}),
	}, nil
}

func (t *tokens) ValidateAccessToken() bool {
	//TODO implement me
	panic("implement me")
}
