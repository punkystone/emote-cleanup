package main

import (
	"go_test/internal"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
	logInstance := os.Getenv("LOG_INSTANCE")
	startDate := os.Getenv("STARTDATE")
	dataDirectory := os.Getenv("DATA_DIRECTORY")
	channel := os.Getenv("CHANNEL")
	err = internal.DownloadLogs(logInstance, startDate, dataDirectory, channel)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
