package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string               `bson:"name" json:"name" binding:"required"`
	Users            []primitive.ObjectID `bson:"users" json:"users"`
	Drifts           []primitive.ObjectID `bson:"drifts" json:"drifts"`
}

func (model *Group) Creating() error {
	if err := model.DefaultModel.Creating(); err != nil {
		return err
	}

	if model.Users == nil {
		model.Users = []primitive.ObjectID{}
	}
	if model.Drifts == nil {
		model.Drifts = []primitive.ObjectID{}
	}

	return nil
}
