package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
	log.Info().Msg(os.Getenv("LOG_INSTANCE"))
	log.Info().Msg(os.Getenv("STARTDATE"))
}
