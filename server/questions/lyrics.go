package questions

import (
	"errors"
	"log"
	"math/rand"

	"github.com/benrhyshoward/hitthespot/model"
	"github.com/google/uuid"
	"github.com/zmb3/spotify"
)

type trackArtist struct {
	TrackName  string
	ArtistName string
}

func generateLyricQuestions(user model.User, noMoreQuestionsChannel chan (struct{})) chan (model.Question) {
	output := make(chan (model.Question))

	go func() {
		log.Print("Starting lyric question generation")
		recentlyPlayed, err := getRecentlyPlayed(user)
		if err != nil {
			log.Println(err.Error())
			return
		}
		topTracks, err := getTopTracks(user)
		if err != nil {
			log.Println(err.Error())
			return
		}

		topArtists, err := getTopArtists(user)
		if err != nil {
			log.Println(err.Error())
			return
		}

		//Collating and deduplicting track and artists
		trackSet := make(map[spotify.ID]trackArtist)
		artistSet := make(map[spotify.ID]string)

		for _, track := range recentlyPlayed {
			if _, found := trackSet[track.Track.ID]; !found {
				trackSet[track.Track.ID] = trackArtist{track.Track.Name, track.Track.Artists[0].Name}
			}
			if _, found := artistSet[track.Track.Artists[0].ID]; !found {
				artistSet[track.Track.Artists[0].ID] = track.Track.Artists[0].Name
			}
		}
		for _, track := range topTracks {
			if _, found := trackSet[track.ID]; !found {
				trackSet[track.ID] = trackArtist{track.Name, track.Artists[0].Name}
			}
			if _, found := artistSet[track.Artists[0].ID]; !found {
				artistSet[track.Artists[0].ID] = track.Artists[0].Name
			}
		}
		for _, artist := range topArtists {
			if _, found := artistSet[artist.ID]; !found {
				artistSet[artist.ID] = artist.Name
			}
		}

		tracks := []trackArtist{}
		for _, v := range trackSet {
			tracks = append(tracks, v)
		}

		artists := []string{}
		for _, v := range artistSet {
			artists = append(artists, v)
		}

		if len(artists) < 4 {
			//if the user hasn't listened to enough different artists to give options, then dont generate any lyric questions
			close(output)
			return
		}

		rand.Shuffle(len(tracks), func(i, j int) {
			tracks[i], tracks[j] = tracks[j], tracks[i]
		})

		for _, track := range tracks {

			artistName := track.ArtistName
			lyrics, err := searchForLyricsByName(track.TrackName, artistName)
			if err != nil {
				continue
			}

			falseOptions, err := getRandomElementsExcluding(artists, 3, track.ArtistName)
			if err != nil {
				continue
			}
			allOptions := append(falseOptions, artistName)
			rand.Shuffle(len(allOptions), func(i, j int) {
				allOptions[i], allOptions[j] = allOptions[j], allOptions[i]
			})

			question := model.Question{
				Id:          uuid.New().String(),
				Type:        model.MultipleChoice,
				Description: "Identify the artist from the lyrics",
				Content:     lyrics,
				Options:     allOptions,
				Answer: model.Answer{
					Value:     artistName,
					ExtraInfo: "From the track '" + track.TrackName + "'",
				},
				Guesses: []model.Guess{},
			}
			select {
			case output <- question:
				log.Print("Sending lyric question to channel")
			case <-noMoreQuestionsChannel:
				log.Print("Stopping lyric question generation")
				return
			}
		}
		close(output)
	}()
	return output
}

func getRandomElementsExcluding(stringSlice []string, n int, excluding string) ([]string, error) {
	var filtered []string
	for _, str := range stringSlice {
		if str != excluding {
			filtered = append(filtered, str)
		}
	}
	if len(filtered) < n {
		return []string{}, errors.New("Not enough elements in array")
	}
	rand.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})
	return filtered[:n], nil
}
