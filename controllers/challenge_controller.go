package controllers

import (
	"net/http"

	"github.com/drift-org/backend/helpers"
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
		Points      int     `json:"points" binding:"required,min=0"`
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
	location.Type = "Point"

	// If address information is provided, include it in our Challenge model.
	if body.Address != "" {
		challenge.Address = body.Address

		// If lat/long was not provided, geocode the address and retrieve the appropriate coordinates.
		if body.Latitude == 0 && body.Longitude == 0 {
			if latitude, longitude, err := helpers.GeocodeAddress(body.Address); err == nil {
				body.Latitude, body.Longitude = latitude, longitude
			} else {
				context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Address or lat/long coordinates are invalid."})
				return
			}
		}
	}

	// If lat/long information is provided/has been geocoded, validate them and include in Challenge model.
	if body.Latitude != 0 && body.Longitude != 0 {
		if helpers.ValidateCoordinates(body.Latitude, body.Longitude) {
			location.Coordinates = []float64{body.Longitude, body.Latitude}
			challenge.Location = &location
		} else {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Lat/long coordinates are invalid."})
			return
		}
	}

	challengeCollection := mgm.Coll(&challenge)

	query := bson.M{"taskName": challenge.TaskName}

	// If location info is included in our challenge, validate duplicates based on this field as well.
	if challenge.Location != nil {
		query["location"] = challenge.Location
	}

	// Check if a duplicate exists - if so, throw an error.
	if err := challengeCollection.First(query, &challenge); err == nil {
		errMsg := "Challenge with same task name already exists."
		if query["location"] != nil {
			errMsg = "Challenge with same task name and location already exists."
		}
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	if err := challengeCollection.Create(&challenge); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Success", "challenge": challenge})
}
