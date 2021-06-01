package models

import (
	"github.com/kamva/mgm/v3"
)

type Challenge struct {
	mgm.DefaultModel `bson:",inline"`
	Longitude        int    `bson:"longitude" json:"longitude"`
	Latitude         int    `bson:"latitude" json:"latitude"`
	Address          string `bson:"address" json:"address"`
	TaskName         string `bson:"taskName" json:"taskName" binding:"required"`
	Description      string `bson:"description" json:"description" binding:"required"`
	Points           int    `bson:"points" json:"points" binding:"required"`
}
