package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct{
	UserID uuid.UUID `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uuid.UUID, role string) (string, error){
	secret :=[]byte(os.Getenv("JWT_SECRET"))

	expirationTime := time.Now().Add(72 * time.Hour)

	claims:= &Claims{
		UserID: userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}