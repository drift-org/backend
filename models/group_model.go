package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	mgm.DefaultModel `bson:",inline"`
	Users            []primitive.ObjectID `bson:"users"`
	Drift            []primitive.ObjectID `bson:"drift"`
}

func (model *Group) Creating() error {
	if err := model.DefaultModel.Creating(); err != nil {
		return err
	}

	if model.Users == nil {
		model.Users = []primitive.ObjectID{}
	}
	if model.Drift == nil {
		model.Drift = []primitive.ObjectID{}
	}

	return nil
}
