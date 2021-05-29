package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userId string, JWTSecret string) (string, error) {

	// Generate claims object that contains user's ObjectID and expires in 1 week.
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userId,
		ExpiresAt: time.Now().Add((time.Hour * 24) * 7).Unix(),
	})

	// Sign the claim and return the generated token string.
	token, err := claims.SignedString([]byte(JWTSecret))
	return token, err
}
