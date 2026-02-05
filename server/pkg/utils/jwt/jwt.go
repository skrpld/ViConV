package jwt

import (
	"sync/atomic"
	"time"
	"viconv/pkg/consts/errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	RefreshTokenExpiryTime time.Duration `env:"REFRESH_TOKEN_EXPIRY_TIME" mapstructure:"REFRESH_TOKEN_EXPIRY_TIME" env-default:"1h"`
	AccessTokenExpiryTime  time.Duration `env:"ACCESS_TOKEN_EXPIRY_TIME" mapstructure:"ACCESS_TOKEN_EXPIRY" env-default:"2h"`
	IssuedAt               string        `env:"ISSUED_AT" mapstructure:"ISSUED_AT" env-default:"viconv"`
}

var currentConfig atomic.Pointer[JWTConfig]

func UpdateJWTConfig(newCfg *JWTConfig) {
	currentConfig.Store(newCfg)
}

func NewAccessToken(id, secret string) (string, error) {
	cfg := currentConfig.Load()
	if cfg == nil {
		return "", errors.ErrInternalServer
	}

	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims = jwt.MapClaims{
		"iss": cfg.IssuedAt,
		"sub": id,
		"exp": time.Now().Add(cfg.AccessTokenExpiryTime).Unix(),
		"iat": time.Now().Unix(),
	}
	return token.SignedString([]byte(secret))
}

func NewRefreshToken(email, secret string) (string, time.Duration, error) {
	cfg := currentConfig.Load()
	if cfg == nil {
		return "", 0, errors.ErrInternalServer
	}

	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims = jwt.MapClaims{
		"iss": cfg.IssuedAt,
		"sub": email,
		"exp": time.Now().Add(cfg.RefreshTokenExpiryTime).Unix(),
		"iat": time.Now().Unix(),
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return tokenString, cfg.RefreshTokenExpiryTime, nil
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
