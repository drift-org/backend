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

// Location (lat/long) and Address information are optional.
// Case 1: Location (provided), Address (provided)
//      - Create object as usual.
// Case 2: Location (provided), Address (not provided)
//      - Create object with Location, but no Address.
// Case 3: Location (not provided), Address (provided)
//      - Geocode address and retrieve lat/long, create
//        object with Location & Address.
// Case 4: Location (not provided), Address (not provided)
//      - Create object as usual, without both information.
//        Validate duplicates based only on name.
type Challenge struct {
	mgm.DefaultModel `bson:",inline"`
	Location         *Location `bson:"location,omitempty" json:"location,omitempty"`
	Address          string    `bson:"address" json:"address"`
	Name             string    `bson:"name" json:"name" binding:"required"`
	Description      string    `bson:"description" json:"description" binding:"required"`
	Points           int       `bson:"points" json:"points" binding:"required,min=0"`
}

func (model *Challenge) Saving() error {
	if err := model.DefaultModel.Saving(); err != nil {
		return err
	}

	if model.Location != nil {
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
