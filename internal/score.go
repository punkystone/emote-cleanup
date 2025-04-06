package internal

import (
	"math"
	"time"
)

const countWeight = 0.3
const lastUsedWeight = 0.7

func getMaxCount(emotes map[string]*Emote) int {
	m := 0
	for _, emote := range emotes {
		if emote.Count > m {
			m = emote.Count
		}
	}
	return m
}

func getMinCount(emotes map[string]*Emote) int {
	m := math.MaxInt32
	for _, emote := range emotes {
		if emote.Count < m {
			m = emote.Count
		}
	}
	return m
}

func getMinLastUsedSeconds(emotes map[string]*Emote) float64 {
	m := math.MaxFloat64
	for _, emote := range emotes {
		if emote.LastUsed != nil && time.Since(*emote.LastUsed).Seconds() < float64(m) {
			m = time.Since(*emote.LastUsed).Seconds()
		}
	}
	return m
}

func getMaxLastUsedSeconds(emotes map[string]*Emote) float64 {
	m := 0.0
	for _, emote := range emotes {
		if emote.LastUsed != nil && time.Since(*emote.LastUsed).Seconds() > m {
			m = time.Since(*emote.LastUsed).Seconds()
		}
	}
	return m
}

func CalculateScores(emotes map[string]*Emote) {
	maxCount := getMaxCount(emotes)
	minCount := getMinCount(emotes)
	maxLastUsedSeconds := getMaxLastUsedSeconds(emotes)
	minLastUsedSeconds := getMinLastUsedSeconds(emotes)
	for _, emote := range emotes {
		var normalizedCount float64
		if maxCount == minCount {
			normalizedCount = 0.5
		} else {
			normalizedCount = float64(emote.Count-minCount) / float64(maxCount-minCount)
		}
		normalizedLastUsed := 0.0
		if emote.LastUsed != nil {
			lastUsedSeconds := time.Since(*emote.LastUsed).Seconds()
			if maxLastUsedSeconds == minLastUsedSeconds {
				normalizedLastUsed = 0.5
			} else {
				normalizedLastUsed = 1 - (lastUsedSeconds-minLastUsedSeconds)/(maxLastUsedSeconds-minLastUsedSeconds)
			}
		}
		emote.Score = countWeight*normalizedCount + lastUsedWeight*normalizedLastUsed
	}
}
