package routes

import (
	"github.com/drift-org/backend/controllers"
	"github.com/drift-org/backend/middleware"
	"github.com/gin-gonic/gin"
)

var (
	prizeController controllers.PrizeController = controllers.NewPrizeController()
)

func prizeRoute(g *gin.RouterGroup) {

	/*
	   In the future, when account roles are added and ADMIN status
	   can be validated, additional auth middleware will be inserted in.
	*/
	g.POST("/", middleware.VerifyAuthenticated(), prizeController.Create)
}
