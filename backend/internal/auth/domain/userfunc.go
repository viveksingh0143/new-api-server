package domain

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vamika-digital/wms-api-server/config"
)

func (u *User) GenerateAccessToken(cfg config.Config) (string, error) {
	expirationTime := time.Now().Add(time.Second * time.Duration(cfg.Auth.ExpiryDuration)).Unix()
	claims := &jwt.StandardClaims{
		Subject:   u.Email,
		ExpiresAt: expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Auth.SecretKey))
}

func (u *User) GenerateRefreshToken(cfg config.Config, expireLong bool) (string, error) {
	var expirationTime int64
	if expireLong {
		expirationTime = time.Now().Add(time.Hour * 24 * time.Duration(cfg.Auth.ExpiryLongDuration)).Unix()
	} else {
		expirationTime = time.Now().Add(time.Second * time.Duration(cfg.Auth.ExpiryDuration)).Unix()
	}

	claims := &jwt.StandardClaims{
		Subject:   u.Email,
		ExpiresAt: expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Auth.SecretKey))
}
