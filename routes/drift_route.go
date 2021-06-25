package routes

import (
	"github.com/drift-org/backend/controllers"
	"github.com/drift-org/backend/middleware"
	"github.com/gin-gonic/gin"
)

var (
	driftController controllers.DriftController = controllers.NewDriftController()
)

func driftRoute(g *gin.RouterGroup) {
	g.POST("/", middleware.VerifyAuthenticated(), driftController.Create)
}
