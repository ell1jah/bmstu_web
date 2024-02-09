package jwt

import (
	"time"

	"github.com/ell1jah/bmstu_web/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type UserClaims struct {
	ID    uint64 `json:"id"`
	Login string `json:"login"`
}

func FromModelUsertoUserClaims(user *model.User) *UserClaims {
	return &UserClaims{
		ID:    user.ID,
		Login: user.Login,
	}
}

type Claims struct {
	User UserClaims `json:"user"`
	jwt.RegisteredClaims
}

type jwtSessionsManager struct {
	key    []byte
	method jwt.SigningMethod
}

func NewJWTSessionsManager(key []byte, method jwt.SigningMethod) *jwtSessionsManager {
	return &jwtSessionsManager{
		key:    key,
		method: method,
	}
}

func (jsm *jwtSessionsManager) CreateSession(user *UserClaims) (string, error) {
	claims := Claims{
		*user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jsm.method, claims)

	tokenString, err := token.SignedString(jsm.key)
	if err != nil {
		return "", errors.Wrap(err, "jwt error")
	}

	return tokenString, nil
}
