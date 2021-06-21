package routes

import (
	"github.com/drift-org/backend/controllers"
	"github.com/gin-gonic/gin"
)

var (
	prizeController controllers.PrizeController = controllers.NewPrizeController()
)

func prizeRoute(g *gin.RouterGroup) {
	g.POST("/create", prizeController.Create)
}
