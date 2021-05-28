package helpers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type JWTService interface {
	GenerateToken(user_id string) (string, error)
	// ValidateToken(tokenString string) (*jwt.Token, error)
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

	
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:    user_id,
		ExpiresAt: time.Now().Add((time.Hour * 24) * 7).Unix(),
	})
	token, err := claims.SignedString([]byte(jwtSrv.secretKey))
	return token, err
}

// func (jwtSrv *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
// 	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// Signing method validation
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		}
// 		// Return the secret signing key
// 		return []byte(jwtSrv.secretKey), nil
// 	})
// }
