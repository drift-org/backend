package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name" json:"name" binding:"required"`
	Age              int    `bson:"age" json:"age" binding:"required"`
	EmailAddress     string `bson:"emailAddress" json:"emailAddress" binding:"required"`
	Password         string `bson:"passsword" json:"passsword" binding:"required"`
	// University       string `bson:"university" binding:"required"`
	Points 		 int `bson:"points" json:"points"`
	Username	 string `bson:"username" json:"username" binding:"required"`
	Friends		 []primitive.ObjectID `bson:"friends" json:"friends"`
	FriendRequests	 []primitive.ObjectID `bson:"friendRequests" json:"friendRequests"`
	PastDrifts	 []primitive.ObjectID `bson:"pastDrifts" json:"pastDrifts"`
	CurrentDrift	 primitive.ObjectID `bson:"currentDrift" json:"currentDrift"`
}

func (model *User) Creating() error {
	if err := model.DefaultModel.Creating(); err != nil {
		return err
	}

	if model.Friends == nil {
		model.Friends = []primitive.ObjectID{}
	}
	if model.FriendRequests == nil {
		model.FriendRequests = []primitive.ObjectID{}
	}
	if model.PastDrifts == nil {
		model.PastDrifts = []primitive.ObjectID{}
	}

	return nil
}