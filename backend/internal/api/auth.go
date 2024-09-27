package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const USER_ID_KEY = "user_id"

func DiscordAuth(ctx *gin.Context) {

	// Check the Authorization header for the token
	authToken := ctx.GetHeader("Authorization")

	// Make a GET request to the Discord API to get the user's info
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	// Check the response
	if resp.StatusCode != 200 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}
	defer resp.Body.Close()

	// Get the user's ID
	body := DiscordAuthBody{}
	data, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Set(USER_ID_KEY, body.UserId)
	ctx.Next()
}

type DiscordAuthBody struct {
	UserId string `json:"id"`
}
