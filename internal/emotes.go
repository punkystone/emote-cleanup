package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

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
			} `json:"connections"`
			EmoteSets []struct {
				ID     string `json:"id"`
				Emotes []struct {
					Name string `json:"name"`
				} `json:"emotes"`
			} `json:"emote_sets"`
		} `json:"userByConnection"`
	} `json:"data"`
}

func GetEmotes(userID string) ([]string, error) {
	seventTVRequest := SevenTVRequest{
		Query: "query($id: String!) { userByConnection(platform: TWITCH, id: $id) { connections { emote_set_id } emote_sets { id emotes { name } } } }",
		Variables: struct {
			ID string `json:"id"`
		}{
			ID: userID,
		},
	}
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(seventTVRequest)
	if err != nil {
		return []string{}, fmt.Errorf("error encoding 7TV request: %w", err)
	}
	response, err := http.Post("https://7tv.io/v3/gql", "application/json", buffer)
	if err != nil {
		return []string{}, fmt.Errorf("error fetching 7TV emotes: %w", err)
	}
	defer response.Body.Close()
	seventTVResponse := SevenTVResponse{}
	err = json.NewDecoder(response.Body).Decode(&seventTVResponse)
	if err != nil {
		return []string{}, fmt.Errorf("error decoding 7TV response: %w", err)
	}
	connections := seventTVResponse.Data.UserByConnection.Connections
	if len(connections) == 0 {
		return []string{}, errors.New("no connections found for user")
	}
	emoteSetID := connections[0].EmoteSetID
	emotes := []string{}
	for _, emoteSet := range seventTVResponse.Data.UserByConnection.EmoteSets {
		if emoteSet.ID == emoteSetID {
			for _, emote := range emoteSet.Emotes {
				emotes = append(emotes, emote.Name)
			}
			break
		}
	}
	return emotes, nil
}
