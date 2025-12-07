package token

import (
	"room-planner/app"
	"room-planner/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey: secretKey}
}

type Token string

func (m *JWTMaker) GenerateToken(user *model.User, ttl time.Duration) (*string, *UserClaims, error) {
	userClaims, err := NewUserClaims(user.Id, user.Username, user.Email, ttl)
	if err != nil {
		return nil, nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenStr, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return nil, nil, err
	}

	return &tokenStr, userClaims, nil
}

func (m *JWTMaker) VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, app.ErrInvalidTokenString
		}

		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, app.ErrInvalidTokenClaims
	}

	return claims, nil
}
