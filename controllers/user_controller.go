package controllers

import (
	"net/http"

	"github.com/drift-org/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
)

type UserController interface {
	Get(context *gin.Context)
}

type userController struct {
}

func NewUserController() UserController {
	return &userController{}
}

func (ctrl *userController) Get(context *gin.Context) {
	// ID refers to MongoDB user object ID. Specificity refers to the verbocity of the returned friends array
	// 0 returns ONLY friend user IDs (default), 1 returns friend user IDs AND usernames
	type IGet struct {
		ID          string `form:"id" binding:"required"`
		Specificity int    `form:"specificity" binding:"eq=0|eq=1"`
	}
	var friendUsernames []string
	var query IGet

	if err := context.ShouldBindQuery(&query); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{}
	coll := mgm.Coll(user)
	err := coll.FindByID(query.ID, user)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Couldn't find user."})
		return
	}

	if query.Specificity == 1 {
		for i := 0; i < len(user.Friends); i++ {
			friend := &models.User{}
			err := coll.FindByID(user.Friends[i], friend)
			if err != nil {
				context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't find friend."})
				return
			}
			friendUsernames = append(friendUsernames, friend.Username)
		}
		context.JSON(http.StatusOK, gin.H{"message": "Success", "user": user, "friendUsernames": friendUsernames})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Success", "user": user})
}
