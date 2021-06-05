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

/*
This function validates a given token string, and if successful, returns a string representation of the ObjectID
of the user who created the token. If unsuccessful, the string is empty, and the error is returned instead.
*/
func ValidateToken(tokenString string) (string, error) {
	// Convert the token string to a JWT token, and verify that it was created with our secret key.
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, (func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtServiceCreds.secretKey), nil
	}))
	if err != nil {
		return "", err
	}
	// Type assertion so that we can extract specific fields from our claims object.
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if ok && token.Valid {
		return claims.Subject, nil
	}
	return "", err

}
