package main

import (
	"go_test/internal"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
	emotesCount := map[string]*internal.Emote{}
	userID := os.Getenv("USER_ID")
	emotes, err := internal.GetEmotes(userID)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	for _, emote := range emotes {
		emotesCount[emote.Name] = &internal.Emote{
			Count:    0,
			LastUsed: nil,
			Added:    emote.Added,
		}
	}
	dataDirectory := os.Getenv("DATA_DIRECTORY")
	logFiles, err := os.ReadDir(dataDirectory)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	for _, logFile := range logFiles {
		err = internal.ScanLogFile(dataDirectory, logFile.Name(), emotesCount)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}
	renderFile := os.Getenv("RENDERFILE")
	internal.Render(emotesCount, renderFile)
}
