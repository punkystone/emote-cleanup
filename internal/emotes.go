package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Emote struct {
	ID       string
	Count    int
	LastUsed *time.Time
	Added    time.Time
}
type SevenTVEmote struct {
	ID    string
	Name  string
	Added time.Time
}

type SevenTVRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ID string `json:"id"`
	} `json:"variables"`
}

type SevenTVResponse struct {
	Data struct {
		UserByConnection struct {
			Connections []struct {
				EmoteSetID string `json:"emote_set_id"`
				Platform   string `json:"platform"`
			} `json:"connections"`
			EmoteSets []struct {
				ID     string `json:"id"`
				Emotes []struct {
					ID        string `json:"id"`
					Name      string `json:"name"`
					Timestamp string `json:"timestamp"`
				} `json:"emotes"`
			} `json:"emote_sets"`
		} `json:"userByConnection"`
	} `json:"data"`
}

func GetEmotes(userID string) ([]SevenTVEmote, error) {
	seventTVRequest := SevenTVRequest{
		Query: "query($id: String!) { userByConnection(platform: TWITCH, id: $id) { connections { emote_set_id platform } emote_sets { id emotes { id name timestamp } } } }",
		Variables: struct {
			ID string `json:"id"`
		}{
			ID: userID,
		},
	}
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(seventTVRequest)
	if err != nil {
		return []SevenTVEmote{}, fmt.Errorf("error encoding 7TV request: %w", err)
	}
	response, err := http.Post("https://7tv.io/v3/gql", "application/json", buffer)
	if err != nil {
		return []SevenTVEmote{}, fmt.Errorf("error fetching 7TV emotes: %w", err)
	}
	defer response.Body.Close()
	seventTVResponse := SevenTVResponse{}
	err = json.NewDecoder(response.Body).Decode(&seventTVResponse)
	if err != nil {
		return []SevenTVEmote{}, fmt.Errorf("error decoding 7TV response: %w", err)
	}
	connections := seventTVResponse.Data.UserByConnection.Connections
	if len(connections) == 0 {
		return []SevenTVEmote{}, errors.New("no connections found for user")
	}
	emoteSetID := ""
	for _, connection := range connections {
		if connection.Platform == "TWITCH" {
			emoteSetID = connection.EmoteSetID
			break
		}
	}
	if emoteSetID == "" {
		return []SevenTVEmote{}, errors.New("no Twitch connection found for user")
	}
	emotes := []SevenTVEmote{}
	for _, emoteSet := range seventTVResponse.Data.UserByConnection.EmoteSets {
		if emoteSet.ID == emoteSetID {
			for _, emote := range emoteSet.Emotes {
				added, err := time.Parse(time.RFC3339, emote.Timestamp)
				if err != nil {
					return []SevenTVEmote{}, fmt.Errorf("error parsing timestamp for emote %s: %w", emote.Name, err)
				}
				emotes = append(emotes, SevenTVEmote{
					Name:  emote.Name,
					Added: added,
					ID:    emote.ID,
				})
			}
			break
		}
	}
	return emotes, nil
}
