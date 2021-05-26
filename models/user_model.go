package models

import "github.com/kamva/mgm/v3"

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name" binding:"required"`
	Age              int    `bson:"age" binding:"required"`
	EmailAddress     string `bson:"email_address" binding:"required"`
	Password         string `bson:"passsword" binding:"required"`
	// University       string `bson:"university" binding:"required"`
	Points           int    `bson:"points"` 
}
