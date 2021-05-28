package models

import (
       "github.com/kamva/mgm/v3"
       "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name" json:"name" binding:"required"`
	Age              int    `bson:"age" json:"age" binding:"required"`
	EmailAddress     string `bson:"email_address" json:"email_address" binding:"required"`
	Password         string `bson:"passsword" json:"passsword" binding:"required"`
	// University       string `bson:"university" binding:"required"`
	Points 		 int `bson:"points" json:"points"`
	Username	 string `bson:"username" json:"username" binding:"required"`
	Friends		 [0]primitive.ObjectID `bson:"friends" json:"friends"`
	FriendRequests	 [0]primitive.ObjectID `bson:"friend_requests" json:"friend_requests"`
	PastDrifts	 [0]primitive.ObjectID `bson:"past_drifts" json:"past_drifts"`
	CurrentDrift	 primitive.ObjectID `bson:"current_drift" json:"current_drift"`
}
