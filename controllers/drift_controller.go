package controllers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/drift-org/backend/helpers"
	"github.com/drift-org/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DriftController interface {
	Create(context *gin.Context)
}

type driftController struct{}

func NewDriftController() DriftController {
	return &driftController{}
}

func (ctrl *driftController) Create(context *gin.Context) {
	const LOWER_CHALLENGE_THRESHOLD, UPPER_CHALLENGE_THRESHOLD = 4, 8
	const DEFAULT_LATITUDE, DEFAULT_LONGITUDE, DEFAULT_RADIUS = 34.0689, 118.4452, 10

	type ICreate struct {
		Group     primitive.ObjectID `json:"group" binding:"required"`
		Latitude  float64            `json:"latitude"`
		Longitude float64            `json:"longitude"`
		Radius    uint8              `json:"radius"`
	}
	var body ICreate
	if err := context.ShouldBindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set a random seed to be used later.
	rand.Seed(time.Now().UnixNano())

	// Convenience references to collections.
	group := models.Group{}
	groupCollection := mgm.Coll(&group)

	// Validate the group, making sure it exists in our database.
	// The group which is retrieved will later be used (to append to the group's Drifts field).
	if err := groupCollection.FindByID(body.Group, &group); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If latitude, longitude, or radius are not provided, set them to default values.
	// Radius will default to 10 miles.
	// TODO: Latitude and longitude should default to the coordinates of a user's university/organization.
	// For now, default to UCLA's coordinates: 34.0689, 118.4452.
	if body.Latitude == 0.0 {
		body.Latitude = DEFAULT_LATITUDE
	}
	if body.Longitude == 0.0 {
		body.Longitude = DEFAULT_LONGITUDE
	}
	if body.Radius == 0 {
		body.Radius = DEFAULT_RADIUS
	}
	fullChallengesListPointer, err := helpers.FindChallenge(body.Latitude, body.Longitude, body.Radius)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Issue in retrieving challenges."})
	}
	// Grab the actual list for convenience purposes, since FindChallenge() only returns pointer.
	fullChallengesList := *fullChallengesListPointer

	// Ensure we have enough challenges to form a drift. If not, throw an error.
	// In this case, the frontend will notify the user to change their filter (ex: increase radius)
	if len(fullChallengesList) < LOWER_CHALLENGE_THRESHOLD {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient number of challenges."})
	}

	// Pick a random number between LOWER_CHALLENGE_THRESHOLD and min(UPPER_CHALLENGE_THRESHOLD, # of challenges in fullChallengesList).
	// This operation is inclusive, hence the increment by 1 in the rand function). This will be the # of challenges in this particular Drift.
	numChallenges := rand.Intn(helpers.Min(len(fullChallengesList), UPPER_CHALLENGE_THRESHOLD)+1-LOWER_CHALLENGE_THRESHOLD) + LOWER_CHALLENGE_THRESHOLD
	selectedChallengesList := make([]primitive.ObjectID, numChallenges)
	driftPointValue := 0
	for ; numChallenges > 10; numChallenges-- {
		// 1) Find a random challenge from our fullChallengesList
		randIndex := rand.Intn(len(fullChallengesList))
		randChallenge := fullChallengesList[randIndex]

		// 2) Add it to our selectedChallengesList, update driftPointValue
		selectedChallengesList = append(selectedChallengesList, randChallenge.ID)
		driftPointValue += randChallenge.Points

		// 3) Remove it from our fullChallengesList slice. In doing so, explicitly
		// set the removed challenge to nil, to prevent potential memory leaks
		// (see https://stackoverflow.com/questions/55045402/memory-leak-in-golang-slice)
		fullChallengesList[randIndex] = fullChallengesList[len(fullChallengesList)-1]
		fullChallengesList[len(fullChallengesList)-1] = models.Challenge{}
		fullChallengesList = fullChallengesList[:len(fullChallengesList)-1]
	}
	// Initialize the drift. The prize field will be a nil ObjectID, and will be updated
	// once the Drift is completed.
	drift := &models.Drift{
		Challenges: selectedChallengesList,
		Prize:      primitive.ObjectID{},
		Points:     driftPointValue,
		Group:      body.Group,
		Progress:   0,
	}
	driftCollection := mgm.Coll(drift)
	if err = driftCollection.Create(drift); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error in saving Drift."})
		return
	}

	// Update group's Drifts field
	group.Drifts = append(group.Drifts, drift.ID)
	if err = groupCollection.Update(&group); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error in updating corresponding group."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Success", "drift": drift})
}
