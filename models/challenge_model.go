package models

import (
	"context"
	"log"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Location struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates" binding:"required"`
}

type Challenge struct {
	mgm.DefaultModel `bson:",inline"`
	Location         Location `bson:"location" json:"location"`
	Address          string   `bson:"address" json:"address"`
	TaskName         string   `bson:"taskName" json:"taskName" binding:"required"`
	Description      string   `bson:"description" json:"description" binding:"required"`
	Points           int      `bson:"points" json:"points" binding:"required"`
}

func (model *Challenge) Saving() error {
	if err := model.DefaultModel.Saving(); err != nil {
		return err
	}

	if model.Location.Coordinates != nil {
		// Set the location type to "Point" for when the coordinates
		// field of challenges are modified (could be during create or update).
		model.Location.Type = "Point"
	}

	return nil
}

func CreateChallengeIndex(ctx context.Context) {
	coll := mgm.Coll(&Challenge{})

	indexView := coll.Indexes()
	model := mongo.IndexModel{Keys: bson.M{"location": "2dsphere"}, Options: nil}
	_, err := indexView.CreateOne(ctx, model)
	if err != nil {
		log.Fatal(err)
	}
}
