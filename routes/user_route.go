package routes

import (
	"github.com/drift-org/backend/controllers"
	"github.com/drift-org/backend/middleware"
	"github.com/gin-gonic/gin"
)

var (
	userController controllers.UserController = controllers.NewUserController()
)

func userRoute(g *gin.RouterGroup) {
	g.GET("/", middleware.VerifyAuthenticated(), userController.Get)
}