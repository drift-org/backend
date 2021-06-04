package routes

import (
	"github.com/drift-org/backend/controllers"
	"github.com/gin-gonic/gin"
)

var (
	groupController controllers.GroupController = controllers.NewGroupController()
)

func groupRoute(g *gin.RouterGroup) {
	g.POST("/group", groupController.Create)
}
