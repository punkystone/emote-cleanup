package internal

import (
	"math"
	"time"
)

const decay = 0.0000001

func CalculateScores(emotes map[string]*Emote) {
	for _, emote := range emotes {
		score := 0.0
		for _, lastUsed := range emote.LastUsed {
			age := time.Since(*lastUsed)
			score += math.Exp(-decay * age.Seconds())
		}
		emote.Score = score
	}
}
