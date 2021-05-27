package controllers

import (
	"net/http"

	"github.com/drift-org/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthController interface {
	Register(context *gin.Context)
}

type authController struct{}

func NewAuthController() AuthController {
	return &authController{}
}

func (ctrl *authController) Register(context *gin.Context) {
	var body models.User
	if err := context.ShouldBindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use bcrypt to encrypt the password.
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.MinCost)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	body.Password = string(encryptedPassword)

	coll := mgm.Coll(&body)

	// Ensure that this email hasn't been registered already.
	err = coll.First(bson.M{"email_address": body.EmailAddress}, &body)
	if err == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User with email already exists."})
		return
	}

	err = coll.Create(&body)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Success", "user": body})
}
