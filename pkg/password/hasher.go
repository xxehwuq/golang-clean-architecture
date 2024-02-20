package password

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const hashingCost = bcrypt.DefaultCost

type Hasher interface {
	Hash(password string) (string, error)
	Compare(password, hashedPassword string) bool
}

type hasher struct {
	salt string
}

func NewHasher(salt string) Hasher {
	return &hasher{
		salt: salt,
	}
}

func (h hasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+h.salt), hashingCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	return string(hashedPassword), nil
}

func (h hasher) Compare(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+h.salt)) == nil
}
