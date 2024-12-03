package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CreatePlaylistRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Artists     []string `json:"artists"`
}

type Playlist struct {
	Collaborative bool              `json:"collaborative"`
	Description   string            `json:"description"`
	ExternalURLS  map[string]string `json:"external_urls"`
	Followers     map[string]any    `json:"followers"`
	Href          string            `json:"href"`
	Id            string            `json:"id"`
	Image         map[string]any    `json:"image"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type ArtistResponse struct {
	Artists struct {
		Href  string   `json:"href"`
		Items []Artist `json:"items"`
	} `json:"artists"`
}

type Artist struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Genres       []string `json:"genres"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	Images     []Image `json:"images"`
	Popularity int     `json:"popularity"`
}

type AddTracksRequest struct {
	Uris []string `json:"uris"`
}

type Track struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	DurationMs int    `json:"duration_ms"`
	Popularity int    `json:"popularity"`
	PreviewURL string `json:"preview_url"`
	Album      struct {
		Name         string            `json:"name"`
		ID           string            `json:"id"`
		ReleaseDate  string            `json:"release_date"`
		TotalTracks  int               `json:"total_tracks"`
		Images       []Image           `json:"images"`
		ExternalURLs map[string]string `json:"external_urls"`
	} `json:"album"`
	Artists []Artist `json:"artists"`
	URI     string   `json:"uri"`
}

type TracksResponse struct {
	Tracks []Track `json:"tracks"`
}

func CreatePlaylistHandler(c *gin.Context) {
	value, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		return
	}

	token, ok := value.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token format"})
		return
	}

	userID := fetchUserID(token)

	playlistData := CreatePlaylistRequest{}
	c.BindJSON(&playlistData)

	playlist := createPlaylist(playlistData, userID, token)

	artistsID := fetchArtist(playlistData, token)

	tracksURI := fetchArtistTopTracks(artistsID, token)

	addTracksToPlaylist(tracksURI, playlist, token)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Playlist created",
	})

}

// The create function creates an empty playlist for the user
func createPlaylist(data CreatePlaylistRequest, userID string, accessToken string) *Playlist {
	body, _ := json.Marshal(data)

	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userID)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	response, _ := client.Do(request)

	playlistResponseBody, _ := io.ReadAll(response.Body)

	playlist := Playlist{}
	json.Unmarshal(playlistResponseBody, &playlist)

	return &playlist
}

func fetchArtist(data CreatePlaylistRequest, accessToken string) []string {
	artistsID := make([]string, len(data.Artists))

	for i, artist := range data.Artists {
		query := strings.Replace(artist, " ", "+", -1)

		url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s+&type=artist", query)

		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		client := &http.Client{}
		response, _ := client.Do(request)

		body, _ := io.ReadAll(response.Body)

		artist := ArtistResponse{}
		json.Unmarshal(body, &artist)

		artistsID[i] = artist.Artists.Items[0].ID
	}

	return artistsID
}

func fetchArtistTopTracks(artistsID []string, accessToken string) []string {
	var tracksURI []string

	for _, artist := range artistsID {
		request, _ := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/artists/%s/top-tracks", artist), nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		client := &http.Client{}
		response, _ := client.Do(request)

		body, _ := io.ReadAll(response.Body)

		tracks := TracksResponse{}
		json.Unmarshal(body, &tracks)

		for _, track := range tracks.Tracks {
			tracksURI = append(tracksURI, track.URI)
		}

	}

	return tracksURI
}

func addTracksToPlaylist(tracksURI []string, playlist *Playlist, accessToken string) {
	requestNum := int(math.Ceil(float64(len(tracksURI)) / 100.0))

	for i := 0; i < requestNum; i++ {
		tracks := AddTracksRequest{}

		if i < requestNum-1 {
			tracks.Uris = tracksURI[i*100 : i*100+100]
		} else {
			tracks.Uris = tracksURI[i*100:]
		}

		body, _ := json.Marshal(tracks)

		request, _ := http.NewRequest("POST", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlist.Id), bytes.NewBuffer(body))

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		client := &http.Client{}
		client.Do(request)
	}
}

func fetchUserID(token string) string {
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, _ := client.Do(req)

	b, _ := io.ReadAll(res.Body)

	spotifyUser := SpotifyUser{}
	json.Unmarshal(b, &spotifyUser)

	return spotifyUser.ID
}
