package main

import (
	"github.com/joho/godotenv"
	"github.com/renatopnasc/made2share-api/internal/config"
	"github.com/renatopnasc/made2share-api/internal/router"
)

func main() {
	// Loading env
	godotenv.Load()

	// Initialize Redis DB connection
	config.Init()

	router.Initialize()
}
