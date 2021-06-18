package main

import (
	"fmt"
	"os"

	"github.com/drift-org/backend/helpers"
	"github.com/drift-org/backend/routes"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	helpers.AlertError(err, "The .env file could not be found")

	// Setup mgm default config
	databaseURL := os.Getenv("MONGO_URL")
	dbName := os.Getenv("MONGO_DBNAME")
	err = mgm.SetDefaultConfig(nil, dbName, options.Client().ApplyURI(databaseURL))
	helpers.AlertError(err, "There was an error connecting to the database")
	fmt.Println("Successfully connected to Mongo")

	helpers.SetupIndexes()

	routes.SetupRouter()
}
