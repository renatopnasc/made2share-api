package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SpotifyUser struct {
	Country         string `json:"country"`
	DisplayName     string `json:"display_name"`
	Email           string `json:"email"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  *string `json:"href"`
		Total int     `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height *int   `json:"height"`
		Width  *int   `json:"width"`
	} `json:"images"`
	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}

func MeHandler(c *gin.Context) {

	token, _ := c.Get("token")

	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data from Spotify"})
		return
	}

	defer res.Body.Close()

	var spotifyUser SpotifyUser
	if err := json.NewDecoder(res.Body).Decode(&spotifyUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Spotify response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": spotifyUser,
	})

}
