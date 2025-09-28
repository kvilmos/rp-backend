package util

import (
	"errors"
	"fmt"
	"room-planner/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Id int64 `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(user model.User) (*string, *string, error) {
	userClaims := JWTClaims{
		Id: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 100)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedAccessToken, err := accessToken.SignedString([]byte("d4mqw2lvfivh7fcr32igzf53"))
	if err != nil {
		return nil, nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 100)),
		},
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte("d4mqw2lvfivh7fcr32igzf53"))
	if err != nil {
		return nil, nil, err
	}

	return &signedAccessToken, &signedRefreshToken, nil
}

func ParseJWTSignedAccessToken(signedAccessToken string) (*JWTClaims, error) {
	parsedJWTAccessToken, err := jwt.ParseWithClaims(signedAccessToken, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("d4mqw2lvfivh7fcr32igzf53"), nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	claims, ok := parsedJWTAccessToken.Claims.(*JWTClaims)
	if !ok {

		fmt.Println(err)
		return nil, errors.New("unknown claims type, cannot proceed")
	}

	return claims, nil
}

func IsClaimExpired(claims *JWTClaims) bool {
	currentTime := jwt.NewNumericDate(time.Now())
	return claims.ExpiresAt.Time.Before(currentTime.Time)
}
