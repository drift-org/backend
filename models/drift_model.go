package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Drift struct {
	mgm.DefaultModel `bson:",inline"`
	Challenges       []primitive.ObjectID `bson:"challenges" json:"challenges" binding:"required"`
	Prize            string               `bson:"prize" json:"prize"`
	Points           int                  `bson:"points" json:"points" binding:"required"`
	Group            primitive.ObjectID   `bson:"group" json:"group" binding:"required"`
	Progress         int                  `bson:"progress" json:"progress" binding:"required"`
}
