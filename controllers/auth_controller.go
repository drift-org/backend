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
	jwtService helpers.JWTService
}

func NewAuthController(jwtService helpers.JWTService) AuthController {
	return &authController{jwtService: jwtService}
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

func (ctrl *authController) Login(context *gin.Context) {
	type ILogin struct {
		EmailAddress string `json:"email_address" binding:"required"`
		Password     string `json:"password" binding:"required"`
	}
	var body ILogin
	if err := context.ShouldBindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := &models.User{}
	coll := mgm.Coll(user)
	err := coll.First(bson.M{"email_address": body.EmailAddress}, user)
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
	user_id := user.ID.String()
	token, err := ctrl.jwtService.GenerateToken(user_id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't generate auth token."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Success", "authToken": token})
}
