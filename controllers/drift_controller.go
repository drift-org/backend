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
		Prize     primitive.ObjectID `json:"prize"`
		Group     primitive.ObjectID `json:"group"`
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
	prize := models.Prize{}
	group := models.Group{}
	prizeCollection := mgm.Coll(&prize)
	groupCollection := mgm.Coll(&group)

	// If the prize field is empty or doesn't exist in our database, select a prize randomly.
	// If there are no prizes in the database, return an error.
	if err := prizeCollection.FindByID(body.Prize, &prize); err != nil {
		// { $sample: { size: <positive integer> } }
		// TODO: Add code for randomly selecting a prize.
		// I
	}
	// Validate the group, making sure it exists in our database.
	// The group which is retrieved will later be used (to append to the group's Drifts field).
	if err := groupCollection.FindByID(body.Group, &group); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Determine when is the optimal time to remove prize from DB - 1) after Drift started,
	// which might mean that we "throw away" prizes to drifts that may never finish, or
	// 2) after Drift completed, which would prevent collisions in case two drifts get the same prize.
	// For now, logic for option 1) is implemented.
	prizeCollection.Delete(&prize)

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

	drift := &models.Drift{
		Challenges: selectedChallengesList,
		Prize:      prize.ID,
		Points:     driftPointValue,
		Group:      body.Group,
		Progress:   0,
	}

	// Update group's Drifts field

	context.JSON(http.StatusOK, gin.H{"message": "Success", "drift": drift})
}
