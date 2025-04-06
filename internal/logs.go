package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"time"

	"github.com/rs/zerolog/log"
)

var messagePattern = regexp.MustCompile(`^\[(\d{4}-\d{2}-\d{2}).+?\] #.+?: (.*)$`)

func DownloadLogs(logInstance string, startDate string, dataDirectory string, channel string) error {
	date, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return errors.New("invalid start date format. Expected YYYY-MM-DD")
	}
	if _, err := os.Stat(dataDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(dataDirectory, 0755)
		if err != nil {
			return fmt.Errorf("error creating data directory: %w", err)
		}
	}
	if err := clearDirectory(dataDirectory); err != nil {
		return fmt.Errorf("error clearing data directory: %w", err)
	}
	now := time.Now()
	for date.Before(now) || date.Equal(now) {
		log.Info().Msgf("Downloading logs for %s", date.Format("2006-01-02"))
		if err := downloadLog(logInstance, channel, date, dataDirectory); err != nil {
			return fmt.Errorf("error downloading log for date %s: %w", date.Format("2006-01-02"), err)
		}
		date = date.AddDate(0, 0, 1)
	}
	return nil
}

func downloadLog(logInstance string, channel string, date time.Time, dataDirectory string) error {
	url := fmt.Sprintf("https://%s/channel/%s/%d/%d/%d", logInstance, channel, date.Year(), date.Month(), date.Day())
	fileName := fmt.Sprintf("%s/%d-%02d-%02d.log", dataDirectory, date.Year(), date.Month(), date.Day())
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}
	defer file.Close()

	//nolint:gosec //allow
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading log: %w", err)
	}
	defer response.Body.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("error writing to log file: %w", err)
	}
	return nil
}

func clearDirectory(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error reading data directory: %w", err)
	}
	for _, file := range files {
		os.Remove(filepath.Join(path, file.Name()))
	}
	return nil
}

func ScanLogFile(dataDirectory string, fileName string, emotesCount map[string]*Emote) error {
	log.Info().Msgf("Scanning %s", fileName)
	file, err := os.Open(dataDirectory + "/" + fileName)
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := messagePattern.FindStringSubmatch(line)
		const requiredGroups = 3
		if len(matches) != requiredGroups {
			continue
		}
		message := matches[2]
		date, err := time.Parse("2006-01-02", matches[1])
		if err != nil {
			log.Error().Msgf("Error parsing date: %v", err)
			continue
		}
		for _, word := range strings.Split(message, " ") {
			if emote, ok := emotesCount[word]; ok {
				emote.Count++
				emote.LastUsed = &date
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log: %w", err)
	}
	return nil
}
