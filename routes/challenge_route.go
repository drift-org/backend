package routes

import (
	"github.com/drift-org/backend/controllers"
	"github.com/gin-gonic/gin"
)

var (
	challengeController controllers.ChallengeController = controllers.NewChallengeController()
)

func challengeRoute(g *gin.RouterGroup) {

	/*
	   In the future, when account roles are added and ADMIN status
	   can be validated, additional auth middleware will be inserted in.
	*/
	g.POST("/", challengeController.Create)
}
