package main

import (
	"log"

	"github.com/Yordi-SE/FlightSearch/config"          // Package for loading configuration
	"github.com/Yordi-SE/FlightSearch/delivery/router" // Package for setting up HTTP routes
	"github.com/Yordi-SE/FlightSearch/use_case"        // Package containing business logic and Sabre client
	"github.com/joho/godotenv"                         // Package for loading .env files
)

// main is the entry point of the application
func main() {
	// Load environment variables from .env file
	// This allows sensitive data like API keys to be stored outside of source code
	err := godotenv.Load()
	if err != nil {
		// If loading fails, log the error and terminate the application
		log.Fatal("Error loading .env file", err)
	}

	// Initialize configuration from environment variables or other sources
	Config, err := config.New()
	if err != nil {
		// If config initialization fails, log the error and exit
		log.Fatal("Error loading config", err)
	}

	// Create a new Sabre client instance with configuration details
	// The client will be used to interact with Sabre's API
	FlightClient := use_case.NewSabreClient(
		Config.ClientID,     // Sabre API Client ID
		Config.ClientSecret, // Sabre API Client Secret
		Config.PCC,          // Pseudo City Code for agency identification
		Config.URL,          // Sabre API endpoint URL
	)

	// Attempt to retrieve an authentication token from Sabre
	// This is required before making any API calls
	error := FlightClient.GetToken()
	if error != nil {
		// If token retrieval fails, log the error and terminate
		log.Fatal("Error getting token", error)
	}

	// Initialize and start the HTTP router with the Sabre client
	// This sets up the web server and API endpoints
	router.NewRouter(FlightClient)
}
