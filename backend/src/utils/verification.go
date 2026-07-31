
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateVerificationToken(userID string) (string, error) {

	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func GenerateResetPasswordToken(userID string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "password_reset",
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateResetPasswordToken(tokenStr string) (jwt.MapClaims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if claims["type"] != "password_reset" {
		return nil, errors.New("invalid reset token")
	}

	return claims, nil
}

