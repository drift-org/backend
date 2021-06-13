package controllers

import (
	"net/http"

	"github.com/drift-org/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
)

type UserController interface {
	Info(context *gin.Context)
}

type userController struct {
}

func NewUserController() UserController {
	return &userController{}
}

func (ctrl *userController) Info(context *gin.Context) {
	// Id refers to MongoDB user object ID. Specificty refers to the verbocity of the returned friends array
	// 0 returns ONLY friend user IDs, 1 returns friend user IDs AND usernames
	type ICreate struct {
		Id          string `json:"id"`
		Specificity int    `json:"specificity"`
	}
	var friendUsernames []string
	var body ICreate

	if err := context.ShouldBindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Specificity > 1 || body.Specificity < 0 {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Specificity must be either 0 or 1"})
		return
	}

	user := &models.User{}
	coll := mgm.Coll(user)
	err := coll.FindByID(body.Id, user)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Couldn't find user."})
		return
	}

	if body.Specificity == 1 {
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
