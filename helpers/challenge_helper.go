package helpers

import (
	"github.com/drift-org/backend/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func milesToMeter(miles uint8) float64 {
	return float64(miles) * 1609.344
}

/*
This helper function finds all the challenges that are within
	radius(miles) based on the given location. 

Parameters:
- longitude is the longitude of the given location 
- latitude is the latitude of the given location 
- radius defines the distance(in miles) the challenges have to be within
*/
func FindChallenge(longitude float64, latitude float64, radius uint8) (*[]models.Challenge, error) {
	coll := mgm.Coll(&models.Challenge{})
	results := []models.Challenge{}
	distance := milesToMeter(radius)

	location := models.Location{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}
	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry":    location,
				"$minDistance": 0,
				"$maxDistance": distance,
			},
		},
	}

	err := coll.SimpleFind(&results, filter)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
