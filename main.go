package main

import (
	"github.com/drift-org/backend/database"
	"github.com/drift-org/backend/helpers"
	"github.com/drift-org/backend/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	helpers.AlertError(err, "The .env file could not be found")

	database.ConnectDB()
	routes.SetupRouter()
}