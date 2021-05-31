package helpers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtServiceCreds JWTServiceCreds = NewJWTServiceCreds()
)

type JWTServiceCreds struct {
	secretKey string
	issuer    string
}

func NewJWTServiceCreds() JWTServiceCreds {
	JWTSecret := os.Getenv("JWT_SECRET")
	return JWTServiceCreds{
		secretKey: JWTSecret,
		issuer:    "drift",
	}
}

/*
This function generates an authentication token for a user going through the login process.
*/
func GenerateToken(userId string) (string, error) {

	// Generate claims object that contains user's ObjectID and expires in 1 week.
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userId,
		ExpiresAt: time.Now().Add((time.Hour * 24) * 7).Unix(),
	})

	// Sign the claim and return the generated token string.
	token, err := claims.SignedString([]byte(jwtServiceCreds.secretKey))
	return token, err
}
