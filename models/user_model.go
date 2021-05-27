package models

import "github.com/kamva/mgm/v3"

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name" json:"name" binding:"required"`
	Age              int    `bson:"age" json:"age" binding:"required"`
	EmailAddress     string `bson:"email_address" json:"email_address" binding:"required"`
	Password         string `bson:"passsword" json:"passsword" binding:"required"`
	// University       string `bson:"university" binding:"required"`
	Points int `bson:"points" json:"points"`
}
