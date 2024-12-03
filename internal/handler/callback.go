package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/renatopnasc/made2share-api/internal/config"
)

func CallbackHandler(c *gin.Context) {
	frontendURI := os.Getenv("FRONTEND_URI")

	code := c.Query("code")
	state := c.Query("state")

	if !states[state] {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "state mismatch"})
		return
	}

	delete(states, state)

	accessToken, _ := exchangeCodeForToken(code)

	sessionID, err := createSession(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	c.SetCookie("_HttpSID", sessionID, 3540, "/", "http://localhost:5173", false, true)

	c.Redirect(http.StatusPermanentRedirect, frontendURI)

}

func exchangeCodeForToken(code string) (string, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}
	return "", fmt.Errorf("no token found")
}

func createSession(accessToken string) (string, error) {
	sessionID := uuid.New().String()

	err := config.GetRedisDB().Set(config.Ctx, sessionID, accessToken, time.Hour).Err()

	return sessionID, err
}
