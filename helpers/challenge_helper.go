package helpers

import (
	"github.com/drift-org/backend/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func FindChallenge(longitude float64, latitude float64, radius uint8) (*[]models.Challenge, error) {
	coll := mgm.Coll(&models.Challenge{})
	results := []models.Challenge{}

	location := models.Location{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}
	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry":    location,
				"$minDistance": 0,
				"$maxDistance": radius,
			},
		},
	}

	err := coll.SimpleFind(&results, filter)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
