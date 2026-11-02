package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims mendefinisikan klaim payload di dalam token JWT.
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT membuat JWT token baru yang ditandatangani dengan secret key.
func GenerateJWT(userID uint, email, name, role, secretKey string, expHours int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expHours) * time.Hour)

	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "kafeinarts-erp",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT memvalidasi token JWT dan mengembalikan klaim jika valid.
func ValidateJWT(tokenString, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode penandatanganan tidak valid")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return claims, nil
}
