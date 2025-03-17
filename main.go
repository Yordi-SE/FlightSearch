package main

import (
	"log"

	"github.com/Yordi-SE/FlightSearch/config"
	"github.com/Yordi-SE/FlightSearch/delivery/router"
	"github.com/Yordi-SE/FlightSearch/use_case"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	Config, err := config.New()
	if err != nil {
		log.Fatal("Error loading config", err)
	}

	SabreClient := use_case.NewSabreClient(Config.ClientID, Config.ClientSecret, Config.PCC,Config.URL)
	error := SabreClient.GetToken()
	if error != nil {
		log.Fatal("Error getting token", error)
	}
	router.NewRouter(SabreClient)
}
