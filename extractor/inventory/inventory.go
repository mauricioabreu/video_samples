package inventory

import (
	"encoding/json"
	"io"
	"net/http"
)

type Streaming struct {
	Name     string `json:"name"`
	Playlist string `json:"playlist"`
}

// GetStreams makes HTTP request to an endpoint in order to get
// the inventory of availables videos to extract resources.
func GetStreams(url string) ([]Streaming, error) {
	var streamings []Streaming
	resp, err := http.Get(url) //nolint
	if err != nil {
		return streamings, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return streamings, err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, &streamings); err != nil {
		return streamings, err
	}

	return streamings, nil
}
