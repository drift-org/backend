package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

var (
	authController controllers.AuthController = controllers.NewAuthController()
)

func authRoute(g *gin.RouterGroup) {
	g.POST("/register", authController.Register)
	g.POST("/login", authController.Login)
}
