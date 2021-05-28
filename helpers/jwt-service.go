package helpers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type JWTService interface {
	GenerateToken(user_id string) (string, error)
}

type JWTServiceInfo struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	envError := godotenv.Load()
	AlertError(envError, "The .env file could not be found")
	JWTSecret := os.Getenv("JWT_SECRET")
	return &JWTServiceInfo{
		secretKey: JWTSecret,
		issuer:    "drift",
	}
}

func (jwtSrv *JWTServiceInfo) GenerateToken(user_id string) (string, error) {
	// Generate claims object that contains user's ObjectID and expires in 1 week.
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   user_id,
		ExpiresAt: time.Now().Add((time.Hour * 24) * 7).Unix(),
	})

	// Sign the claim and return the generated token string.
	token, err := claims.SignedString([]byte(jwtSrv.secretKey))
	return token, err
}
