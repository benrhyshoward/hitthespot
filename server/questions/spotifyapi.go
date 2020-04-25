package questions

import (
	"github.com/benrhyshoward/hitthespot/model"
	"github.com/zmb3/spotify"
)

func getTopArtists(user model.User) ([]spotify.FullArtist, error) {
	limit := 50
	timeRange := "long"

	longTermTopArtists, err := user.Client.CurrentUsersTopArtistsOpt(&spotify.Options{
		Timerange: &timeRange,
		Limit:     &limit,
	})
	if err != nil {
		return nil, err
	}
	timeRange = "medium"
	mediumTermTopArtists, err := user.Client.CurrentUsersTopArtistsOpt(&spotify.Options{
		Timerange: &timeRange,
		Limit:     &limit,
	})
	if err != nil {
		return nil, err
	}
	timeRange = "short"
	shortTermTopArtists, err := user.Client.CurrentUsersTopArtistsOpt(&spotify.Options{
		Timerange: &timeRange,
		Limit:     &limit,
	})
	if err != nil {
		return nil, err
	}

	//appending artist lists and deduplicating
	artistSet := make(map[spotify.ID]spotify.FullArtist)

	for _, artist := range longTermTopArtists.Artists {
		if _, found := artistSet[artist.ID]; !found {
			artistSet[artist.ID] = artist
		}
	}
	for _, artist := range mediumTermTopArtists.Artists {
		if _, found := artistSet[artist.ID]; !found {
			artistSet[artist.ID] = artist
		}
	}
	for _, artist := range shortTermTopArtists.Artists {
		if _, found := artistSet[artist.ID]; !found {
			artistSet[artist.ID] = artist
		}
	}

	//getting value set from map
	artists := []spotify.FullArtist{}
	for _, v := range artistSet {
		artists = append(artists, v)
	}

	return artists, nil
}

func getTopTracks(user model.User) ([]spotify.FullTrack, error) {
	limit := 50
	timeRange := "long"
	longTermTopTracks, err := user.Client.CurrentUsersTopTracksOpt(&spotify.Options{
		Timerange: &timeRange,
		Limit:     &limit,
	})
	if err != nil {
		return nil, err
	}
	timeRange = "medium"
	mediumTermTopTracks, err := user.Client.CurrentUsersTopTracksOpt(&spotify.Options{
		Timerange: &timeRange,
		Limit:     &limit,
	})
	if err != nil {
		return nil, err
	}
	timeRange = "short"
	shortTermTopTracks, err := user.Client.CurrentUsersTopTracksOpt(&spotify.Options{
		Timerange: &timeRange,
		Limit:     &limit,
	})
	if err != nil {
		return nil, err
	}

	//appending artist lists and deduplicating
	trackSet := make(map[spotify.ID]spotify.FullTrack)

	for _, track := range longTermTopTracks.Tracks {
		if _, found := trackSet[track.ID]; !found {
			trackSet[track.ID] = track
		}
	}
	for _, track := range mediumTermTopTracks.Tracks {
		if _, found := trackSet[track.ID]; !found {
			trackSet[track.ID] = track
		}
	}
	for _, track := range shortTermTopTracks.Tracks {
		if _, found := trackSet[track.ID]; !found {
			trackSet[track.ID] = track
		}
	}

	//getting value set from map
	tracks := []spotify.FullTrack{}
	for _, v := range trackSet {
		tracks = append(tracks, v)
	}
	return tracks, nil
}

func getRecentlyPlayed(user model.User) ([]spotify.RecentlyPlayedItem, error) {
	limit := 50
	recentlyPlayed, err := user.Client.PlayerRecentlyPlayedOpt(&spotify.RecentlyPlayedOptions{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	return recentlyPlayed, nil

}
