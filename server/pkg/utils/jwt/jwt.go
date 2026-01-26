package jwt

import (
	"time"
	"viconv/pkg/consts/errors"

	"github.com/golang-jwt/jwt/v5"
)

func NewAccessToken(id, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims = jwt.MapClaims{
		"iss": "viconv", //TODO: вынести в real-time config
		"sub": id,
		"exp": time.Now().Add(time.Hour * 2).Unix(), //TODO: вынести в real-time config
		"iat": time.Now().Unix(),
	}
	return token.SignedString([]byte(secret))
}

func NewRefreshToken(email, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims = jwt.MapClaims{
		"iss": "viconv", //TODO: вынести в real-time config
		"sub": email,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), //TODO: вынести в real-time config
		"iat": time.Now().Unix(),
	}
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string, secret string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	if token == nil {
		return nil, errors.ErrInvalidToken
	}

	if !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.ErrInvalidToken
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return nil, errors.ErrInvalidToken
	}
	if time.Now().After(exp.Time) {
		return nil, errors.ErrExpiredToken
	}

	return &claims, nil
}
