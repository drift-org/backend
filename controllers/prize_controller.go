package controllers

import (
	"net/http"

	"github.com/drift-org/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
)

type PrizeController interface {
	Create(context *gin.Context)
}

type prizeController struct{}

func NewPrizeController() PrizeController {
	return &prizeController{}
}

func (ctrl *prizeController) Create(context *gin.Context) {
	var body models.Prize
	if err := context.ShouldBindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	prizeCollection := mgm.Coll(&models.Prize{})
	err := prizeCollection.Create(&body)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Success", "prize": body})
}
