package api

import (
	"keno/internal/db"
	"keno/internal/engine"
	"keno/internal/models"
	"keno/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

// Place your Keno Picks
// @Summary Place your picks for the next Keno game
// @Description Give us your numbers so you can enjoy the number of games you specify. There are some rules:
// @Description - You can only pick numbers between `1` and `80`.
// @Description - You can only pick `1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20, 40` numbers per game.
// @Description - You can only play `1, 2, 3, 4, 5, 10, 20, 50, 100` number of games.
// @Tags picks
// @Accept json
// @Produce json
// @Param picks body PickRequest true "Your picks for the next selected games"
// @Success 200 {object} PickResponse
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
// @Router /api/v1/picks [post]
func PlacePicks(ctx *gin.Context) {

	// Get the Picks from the request
	req := PickRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.WithField("src", "api.PlacePicks").Error("Picks call made with invalid JSON body")
		ctx.JSON(http.StatusBadRequest, ErrInvalidPicks)
		return
	}

	// Validate the picks
	if !req.isValid() {
		log.WithField("src", "api.PlacePicks").Error("Picks call made with invalid values")
		ctx.JSON(http.StatusBadRequest, ErrInvalidPicks)
		return
	}

	// Get the database from the context
	db, ok := ctx.Get(db.DbKey)
	if !ok {
		log.WithField("src", "api.PlacePicks").Error("Database not found in context")
		ctx.JSON(http.StatusInternalServerError, ErrInternalError)
		return
	}

	// Get the game engine from the context
	gameEngine, ok := ctx.Get(engine.EngineKey)
	if !ok {
		log.WithField("src", "api.PlacePicks").Error("Game Engine not found in context")
		ctx.JSON(http.StatusInternalServerError, ErrInternalError)
		return
	}

	// Place the picks
	card, err := models.SubmitCard(
		db.(*gorm.DB),
		req.Picks,
		gameEngine.(*engine.Engine).GetGameNumber(),
		req.NumGames,
		req.PricePerGame,
	)
	if err != nil {
		log.WithField("src", "api.PlacePicks").Error("Error submitting picks")
		ctx.JSON(http.StatusInternalServerError, ErrInternalError)
		return
	}

	// Return the card
	ctx.JSON(http.StatusOK, cardToPickResponse(*card))
}

type PickRequest struct {
	PicksPerGame uint8   `json:"picks_per_game"`
	Picks        []uint8 `json:"picks"`
	PricePerGame uint64  `json:"price_per_game"`
	NumGames     uint8   `json:"number_games"`
}

var (
	ValidPickMin      uint8 = 1
	ValidPickMax      uint8 = 80
	ValidPicksPerGame       = []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20, 40}
	ValidGames              = []uint8{1, 2, 3, 4, 5, 10, 20, 50, 100}
)

func (p PickRequest) isValid() bool {
	// Make sure all picks are in range 1-80
	for _, num := range p.Picks {
		if num < ValidPickMin || num > ValidPickMax {
			return false
		}
	}

	// Make sure the number of picks is valid
	if !utils.Contains(ValidPicksPerGame, p.PicksPerGame) {
		return false
	}

	// Make sure the number of games is valid
	if !utils.Contains(ValidGames, p.NumGames) {
		return false
	}

	// Check if numbers selected match the number of picks
	if uint8(len(p.Picks)) != p.PicksPerGame {
		return false
	}

	return true
}

type PickResponse struct {
	CardId    uint64 `json:"card_id"`
	Selection []int  `json:"selection"`
	StartGame uint64 `json:"start_game_num"`
	LastGame  uint64 `json:"last_game_num"`
}

func cardToPickResponse(card models.Card) PickResponse {
	resp := PickResponse{
		CardId:    card.ID,
		Selection: make([]int, 0),
		StartGame: card.StartGame,
		LastGame:  card.LastGame,
	}

	for _, num := range card.Selection {
		resp.Selection = append(resp.Selection, int(num))
	}
	return resp
}
