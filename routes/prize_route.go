package routes

import (
	"github.com/drift-org/backend/controllers"
	"github.com/gin-gonic/gin"
)

var (
	prizeController controllers.PrizeController = controllers.NewPrizeController()
)

func prizeRoute(g *gin.RouterGroup) {

	/*
	   For now, authentication is not necessary to create a prize.
	   In the future, when account roles are added and ADMIN status
	   can be validated, auth middleware will be inserted in. Until then,
	   this route is largely intended for testing purposes.
	*/
	g.POST("/", prizeController.Create)
}
