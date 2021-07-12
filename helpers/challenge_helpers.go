package helpers

import (
	"github.com/drift-org/backend/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// Default lat/long are given as UCLA's coordinates. Default radius is arbitrarily set to 10.
const DEFAULT_LATITUDE, DEFAULT_LONGITUDE, DEFAULT_RADIUS = 34.0701449, -118.4422936, 10

func milesToMeter(miles uint8) float64 {
	const METERS_PER_MILE = 1609.344
	return float64(miles) * METERS_PER_MILE
}

/*
This helper function finds all the challenges that are within
radius (miles) based on the given location.

Parameters:
- longitude is the longitude of the given location
- latitude is the latitude of the given location
- radius defines the distance (in miles) the challenges have to be within
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

/*
TODO: This helper function geocodes the provided address to a valid set
of lat/long coordinates. Currently, it returns a default response of coordinates,
currently UCLA's location.

Parameters:
- address is the address, in string format, provided by the user
*/
func GeocodeAddress(address string) (float64, float64, error) {
	// Return the default lat/long coordinates for now.
	return DEFAULT_LATITUDE, DEFAULT_LONGITUDE, nil
}

/*
Helper function to determine if lat and long coordinates are valid.

Parameters:
- latitude is the latitude of the given location
- longitude is the longitude of the given location
*/
func ValidateCoordinates(latitude float64, longitude float64) bool {
	// TODO: Possibly add an optional "address" field where we verify that the coordinates & the address point
	// to the same location?
	return (-90 <= latitude && latitude <= 90) && (-180 <= longitude && longitude <= 180)
}
