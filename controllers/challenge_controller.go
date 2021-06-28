package controllers

import (
	"net/http"

	"github.com/drift-org/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type ChallengeController interface {
	Create(context *gin.Context)
}

type challengeController struct{}

func NewChallengeController() ChallengeController {
	return &challengeController{}
}

// TODO: Idea for future, if we want to let companies/orgs create challenges (Business Model #3).
// Default method of adding challenges on frontend is to add by placing a marker on a map
// (allowing us to grab Lat/Long coordinates). We can also offer a (premium) option for users to
// provide just the address, and convert to Lat/Long -- thus allowing batch creation of challenges.
func (ctrl *challengeController) Create(context *gin.Context) {
	type ICreate struct {
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		Address     string  `json:"address"`
		TaskName    string  `json:"taskName" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Points      int     `json:"points" binding:"required,min=1"`
	}
	var body ICreate
	if err := context.ShouldBindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	location := models.Location{}

	challenge := models.Challenge{
		TaskName:    body.TaskName,
		Description: body.Description,
		Points:      body.Points,
	}

	// If lat/long/address information is provided, include it in our Challenge model.
	if body.Latitude != 0 && body.Longitude != 0 && body.Address != "" {
		location.Coordinates = []float64{body.Longitude, body.Latitude}
		challenge.Location = &location
		challenge.Address = body.Address
	}
	challengeCollection := mgm.Coll(&challenge)
	// Check if a challenge with the same address and task exist. If so, throw an error.
	if err := challengeCollection.First(bson.M{"taskName": challenge.TaskName, "address": challenge.Address}, &challenge); err == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Challenge with same address and task name already exists."})
		return
	}

	if err := challengeCollection.Create(&challenge); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Success", "challenge": challenge})
}
