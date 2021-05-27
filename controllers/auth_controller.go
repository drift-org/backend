package controllers

import (
	"github.com/drift-org/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthController interface {
	Register(context *gin.Context)
}

type authController struct{}

func NewAuthController() AuthController {
	return &authController{}
}

func (this *authController) Register(context *gin.Context) {
	type IRegister struct {
		Name         string `json:"name" binding:"required"`
		Age          int    `json:"age" binding:"required"`
		EmailAddress string `json:"email_address" binding:"required"`
		Password     string `json:"passsword" binding:"required"`
	}
	var body IRegister
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

	user := &models.User{
		Name:         body.Name,
		Age:          body.Age,
		EmailAddress: body.EmailAddress,
		Password:     string(encryptedPassword),
	}

	coll := mgm.Coll(user)

	// Ensure that this email hasn't been registered already.
	var existingUsers []models.User
	coll.SimpleFind(&existingUsers, bson.M{"email_address": user.EmailAddress})
	if len(existingUsers) > 0 {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User with email already exists."})
		return
	}

	err = coll.Create(user)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Success", "user": user})
}
