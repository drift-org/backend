package models

import (
	"github.com/kamva/mgm/v3"
)

type Prize struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name" json:"name" binding:"required"`
	Description      string `bson:"description" json:"description" binding:"required"`
}
