package controllers

import (
	"net/http"
	"sort"

	"backend/helpers"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupController interface {
	Create(context *gin.Context)
}

type groupController struct{}

func NewGroupController() GroupController {
	return &groupController{}
}

func (ctrl *groupController) Create(context *gin.Context) {
	type ICreate struct {
		Usernames []string `json:"usernames" binding:"required,min=1,unique"`
	}
	var body ICreate
	if err := context.ShouldBindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sort the usernames. This is so that the SAME userIDs slice is created for different orderings
	// of the usernames passed-in. This is so that we can guarantee that a group isn't created twice
	// with different orderings of users.
	sort.Strings(body.Usernames)

	// Convenience references to collections.
	userColllection := mgm.Coll(&models.User{})
	groupCollection := mgm.Coll(&models.Group{})

	//----------------------------------------------------------------------------

	// Step 1: Find the User models that match the Usernames passed-in.
	// For performance reasons, we use a pipeline filter to find the users with the
	// Usernames.
	users := []models.User{}

	// Map the Usernames to a slice of maps, specifying the usernames
	pipeline := helpers.Map(body.Usernames, func(username interface{}) interface{} {
		return bson.M{"username": username}
	})

	// Find the User models, using the $or filter to find all of the unique usernames.
	err := userColllection.SimpleFind(&users, bson.M{"$or": pipeline})
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make sure that every username was found in the database.
	if len(users) != len(body.Usernames) {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Some usernames were not found."})
		return
	}

	//----------------------------------------------------------------------------------------

	// Step 2: Find the id of each of the User models found.
	// While iterating, make sure that the user who created the group is inside the userIDs.
	userIDs := make([]primitive.ObjectID, len(users))
	bearerUserFound := false

	for i := 0; i < len(userIDs); i++ {
		userIDs[i] = users[i].ID

		if !bearerUserFound && userIDs[i].Hex() == context.MustGet("userID") {
			bearerUserFound = true
		}
	}

	// Make sure that the user who created the group is inside the usernames.
	if !bearerUserFound {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Usernames must contain the bearer user."})
		return
	}

	//----------------------------------------------------------------------------------------

	// Step 3: Create the Group.
	group := &models.Group{
		Users: userIDs,
	}

	// Ensure that this group hasn't been created already. Since the userIDs are sorted and unique, we can
	// search by the users field.
	err = groupCollection.First(bson.M{"users": userIDs}, group)
	if err == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Group already exists."})
		return
	}

	err = groupCollection.Create(group)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Success", "group": group})
}
