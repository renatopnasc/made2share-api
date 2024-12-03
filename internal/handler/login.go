package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const redirectURI = "http://localhost:8080/api/v1/callback"

func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

var states = make(map[string]bool)

func LoginHandler(c *gin.Context) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	scope := "playlist-modify-public playlist-modify-private"

	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}

	states[state] = true

	authURL := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s", clientID, redirectURI, scope, state)

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}
