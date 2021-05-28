package routes

import (
	"github.com/drift-org/backend/controllers"
	"github.com/drift-org/backend/helpers"
	"github.com/gin-gonic/gin"
)

var (
	jwtService   helpers.JWTService   = helpers.NewJWTService()
	authController controllers.AuthController = controllers.NewAuthController(jwtService)
)

func authRoute(g *gin.RouterGroup) {
	g.POST("/register", authController.Register)
	g.POST("/login", authController.Login)
}
