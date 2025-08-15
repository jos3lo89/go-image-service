package lib

import (
	"jos3lo89/go-image-service/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, userName string) (string, error) {
	jwtSecret := []byte(config.AppConfig.JWTSecret)
	claims := JWTClaims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(jwtSecret)
	return ss, err
}
