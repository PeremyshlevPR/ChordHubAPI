package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	*jwt.StandardClaims
}

func IssueToken(userId uint, role, secretKey string, expTimeDuration time.Duration) (string, error) {
	expirationTime := time.Now().Add(expTimeDuration)
	claims := &Claims{
		userId,
		role,
		&jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string, secretKey string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return claims, err
	}
	return claims, nil
}
