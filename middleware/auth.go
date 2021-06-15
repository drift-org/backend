package middleware

import (
	"net/http"
	"strings"

	"backend/helpers"

	"github.com/gin-gonic/gin"
)

/*
Middleware function to be used to protect authenticated routes. Requires authentication token string
to be passed in via Authorization Header in the format: "Bearer <token>"
*/
func VerifyAuthenticated() gin.HandlerFunc {
	return func(context *gin.Context) {

		tokenString := strings.Replace(context.Request.Header.Get("Authorization"), "Bearer ", "", 1)
		userID, err := helpers.ValidateToken(tokenString)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Auth token incorrect or expired."})
			return
		}
		context.Set("userID", userID)
		context.Next()
	}
}
