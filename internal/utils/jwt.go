package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/paaart/kavalife-erp-backend/internal/config"
)

var jwtSecret = []byte(config.ConfigLoader().JWT_SECRET)

func CreateJWT(username string, tokenDuration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(tokenDuration).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenStr string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, true
	}
	return nil, false
}
