package jwttoken

import (
	"errors"
	"fmt"
	"github.com/PengShaw/go-common/ginx/ginx-auth/models"
	"github.com/golang-jwt/jwt"
	"time"
)

var (
	privateKey     string
	expired        int
	TokenHeaderKey string

	ErrTokenFormat    = errors.New("token format is wrong")
	ErrTokenExpired   = errors.New("token expired")
	ErrTokenNoActive  = errors.New("token not active")
	ErrTokenUnHandler = errors.New("couldn't handle this token")
)

func Setup(key string, JwtExpired int, tokenHeaderKey string) {
	privateKey = key
	expired = JwtExpired
	TokenHeaderKey = tokenHeaderKey
}

type Claims struct {
	UserID uint
	jwt.StandardClaims
}

func Generate(user models.AuthUser) (string, error) {
	key := []byte(privateKey)

	claims := Claims{
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expired) * time.Second).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func Verify(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(privateKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenFormat
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNoActive
			}
		}
		return nil, fmt.Errorf("%w: [%w]", ErrTokenUnHandler, err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenUnHandler
}
