package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtManager struct {
	secret        string
	tokenDuration time.Duration
}

func NewJwtManager(secret string, tokenDuration time.Duration) *JwtManager {
	return &JwtManager{
		secret:        secret,
		tokenDuration: tokenDuration,
	}
}

func (jm *JwtManager) GenerateToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(jm.tokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jm.secret))
}

func (jm *JwtManager) ValidateToken(tokenString string) (string, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jm.secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userId, ok := claims["userId"].(string)
		if !ok {
			return "", fmt.Errorf("invalid token claims")
		}
		return userId, nil
	}

	return "", fmt.Errorf("invalid token")
}
