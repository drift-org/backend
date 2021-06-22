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

// mgm.DefaultModel `bson:",inline"`
// Challenges       []primitive.ObjectID `bson:"challenges" json:"challenges" binding:"required"`
// Prize            primitive.ObjectID   `bson:"prize" json:"prize"`
// Points           int                  `bson:"points" json:"points" binding:"required"`
// Group            primitive.ObjectID   `bson:"group" json:"group" binding:"required"`
// Progress         int                  `bson:"progress" json:"progress" binding:"required"`

type DriftController interface {
	Create(context *gin.Context)
}

type driftController struct{}

func NewDriftController() DriftController {
	return &driftController{}
}

func (ctrl *driftController) Create(context *gin.Context) {
	const LOWER_CHALLENGE_THRESHOLD, UPPER_CHALLENGE_THRESHOLD = 5, 8
	const DEFAULT_LATITUDE, DEFAULT_LONGITUDE, DEFAULT_RADIUS = 34.0689, 118.4452, 10

	type ICreate struct {
		Prize     primitive.ObjectID `json:"prize"`
		Points    int                `json:"points" binding:"required,min=1"`
		Group     primitive.ObjectID `json:"group"`
		Longitude float64            `json:"longitude"`
		Latitude  float64            `json:"latitude"`
		Radius    uint8              `json:"radius"`
	}
	var body ICreate
	if err := context.ShouldBindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convenience references to collections.
	prize := models.Prize{}
	group := models.Group{}

	prizeCollection := mgm.Coll(&prize)
	groupCollection := mgm.Coll(&group)
	challengeCollection := mgm.Coll(&models.Challenge{})

	// Validate the group and prize, making sure that both exist in our database.
	// The group which is retrieved will later be used (to append to its Drifts field).
	// The prize which is retrieved will be removed from the DB.
	err := groupCollection.FindByID(body.Group, &group)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = prizeCollection.FindByID(body.Prize, &prize)
	if err != nil {
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
	fullChallengesList, err := helpers.FindChallenge(body.Longitude, body.Latitude, body.Radius)
	var selectedChallengesList []primitive.ObjectID

	// Pick a random number between LOWER_CHALLENGE_THRESHOLD and min(UPPER_CHALLENGE_THRESHOLD, # of challenges in fullChallengesList).
	// This operation is inclusive, hence the increment by 1 in the rand function).This will be the # of challenges in this particular Drift.
	rand.Seed(time.Now().UnixNano())
	var numChallenges = rand.Intn(helpers.Min(len(*fullChallengesList), UPPER_CHALLENGE_THRESHOLD)+1-LOWER_CHALLENGE_THRESHOLD) + LOWER_CHALLENGE_THRESHOLD
	for ; numChallenges > 10; numChallenges-- {
		// 1) Find a random challenge from our fullChallengesList
		// 2) Add it to our selectedChallengesList
		// 3) Remove it from our fullChallengesList pool
		selectedChallengesList = append(selectedChallengesList, primitive.ObjectID{})
	}
	context.JSON(http.StatusOK, gin.H{"message": "Success", "drift": nil})
}
