package controllers

import (
	"net/http"

	"github.com/drift-org/backend/helpers"
	"github.com/drift-org/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthController interface {
	Register(context *gin.Context)
	Login(context *gin.Context)
}

type authController struct {
}

func NewAuthController() AuthController {
	// helpers.AlertError(envError, "The .env file could not be found")
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
	err = coll.First(bson.M{"emailAddress": body.EmailAddress}, &body)
	if err == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User with email already exists."})
		return
	}
	// Ensure that this username hasn't been registered already.
	err = coll.First(bson.M{"username": body.Username}, &body)
	if err == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User with username already exists."})
		return
	}

	err = coll.Create(&body)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Success", "user": body})
}

func (ctrl *authController) Login(context *gin.Context) {
	type ILogin struct {
		EmailAddress string `json:"emailAddress"`
		Username string `json:"username"`
		Password     string `json:"password" binding:"required"`
	}
	var body ILogin
	if err := context.ShouldBindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Determine if user is logging in with email or username
	var loginType string
	var loginValue string
	if body.EmailAddress != "" {
		loginType = "emailAddress"
		loginValue = body.EmailAddress
	} else if body.Username != "" {
		loginType = "username"
		loginValue = body.Username
	} else {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Neither username/email address not provided."})
		return
	}

	user := &models.User{}
	coll := mgm.Coll(user)
	err := coll.First(bson.M{loginType: loginValue}, user)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Couldn't find user."})
		return
	}
	hashedPassword := user.Password
	
	// Compare plaintext password with hashed password stored in DB.
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body.Password)); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect password."})
		return
	}
	userId := user.ID.String()
	token, err := helpers.GenerateToken(userId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't generate auth token."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Success", "authToken": token})
}

