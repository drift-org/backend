package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

var (
	groupController controllers.GroupController = controllers.NewGroupController()
)

func groupRoute(g *gin.RouterGroup) {
	g.POST("/group", middleware.VerifyAuthenticated(), groupController.Create)
}
