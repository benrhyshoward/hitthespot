package questions

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var apiroot string = "https://api.musixmatch.com/ws/1.1/"

var client http.Client

func init() {
	client = http.Client{}
}

//Only extracting the values we want from the response
type TrackSearchResponse struct {
	Message struct {
		Header struct {
			StatusCode int `json:"status_code"`
			Available  int `json:"available"`
		} `json:"header"`
		Body struct {
			TrackList []struct {
				Track struct {
					TrackID int `json:"track_id"`
				} `json:"track"`
			} `json:"track_list"`
		} `json:"body"`
	} `json:"message"`
}

type GetLyricsResponse struct {
	Message struct {
		Header struct {
			StatusCode int `json:"status_code"`
		} `json:"header"`
		Body struct {
			Lyrics struct {
				LyricsBody string `json:"lyrics_body"`
			} `json:"lyrics"`
		} `json:"body"`
	} `json:"message"`
}

func getApiKey() string {
	return os.Getenv("MUSIXMATCH_API_KEY")
}

func searchForTrackByName(trackName string, artistName string) (string, error) {
	request, err := http.NewRequest("GET", apiroot+"track.search", nil)
	if err != nil {
		return "", err
	}

	q := request.URL.Query()
	q.Add("format", "json")
	q.Add("q_track", trackName)
	q.Add("q_artist", artistName)
	q.Add("f_has_lyrics", strconv.Itoa(1))
	q.Add("apikey", getApiKey())
	request.URL.RawQuery = q.Encode()

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseObject TrackSearchResponse
	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&responseObject)

	if responseObject.Message.Header.StatusCode == 200 && responseObject.Message.Header.Available > 0 {
		return strconv.Itoa(responseObject.Message.Body.TrackList[0].Track.TrackID), nil
	}

	return "", errors.New("Track not found")
}

func searchForLyricsById(trackID string) (string, error) {
	request, err := http.NewRequest("GET", apiroot+"track.lyrics.get", nil)
	if err != nil {
		return "", err
	}

	q := request.URL.Query()
	q.Add("format", "json")
	q.Add("track_id", trackID)
	q.Add("apikey", getApiKey())
	request.URL.RawQuery = q.Encode()

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseObject GetLyricsResponse
	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&responseObject)

	if responseObject.Message.Header.StatusCode == 200 && responseObject.Message.Body.Lyrics.LyricsBody != "" {

		return formatLyrics(responseObject.Message.Body.Lyrics.LyricsBody), nil
	}
	return "", errors.New("Lyrics not found")
}

func formatLyrics(lyrics string) string {

	lines := strings.Split(lyrics, "\n")
	lines = lines[:len(lines)-4] //last three lines are are added by the api

	var filteredLines []string
	for _, line := range lines {
		if line != "" && //ignore empty lines
			!(strings.HasPrefix(line, "[") && strings.HasPrefix(line, "]")) { //ignore lines indicating song structure
			filteredLines = append(filteredLines, line)
		}
	}

	if len(filteredLines) > 4 {
		filteredLines = filteredLines[:4]
	}
	return strings.Join(filteredLines, "\n")
}

func searchForLyricsByName(trackName string, artistName string) (string, error) {
	trackID, err := searchForTrackByName(trackName, artistName)
	if err != nil {
		return "", err
	}
	lyrics, err := searchForLyricsById(trackID)
	if err != nil {
		return "", err
	}
	return lyrics, nil
}
