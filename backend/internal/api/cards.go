package api

import (
	"keno/internal/db"
	"keno/internal/engine"
	"keno/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

// Check Card
// @Summary Check your card to see if you won
// @Description Check your card to see if you won and claim your wins.
// @Tags cards
// @param card_id path int true "Card ID"
// @Produce json
// @Success 200 {object} CheckCardResponse
// @Failure 400 {object} APIError
// @Failure 404 {object} APIError
// @Failure 500 {object} APIError
// @Router /api/v1/check/{card_id} [get]
func CheckCard(ctx *gin.Context) {
	// Get Card Id from URL
	cardIdStr := ctx.Param("card_id")
	if cardIdStr == "" {
		ctx.JSON(400, ErrInvalidCard)
		return
	}

	// Convert to uint64
	cardId, err := strconv.ParseUint(cardIdStr, 10, 64)
	if err != nil {
		ctx.JSON(400, ErrInvalidCard)
		return
	}

	// Get the database from the context
	db, ok := ctx.Get(db.DbKey)
	if !ok {
		ctx.JSON(500, ErrInternalError)
		return
	}

	// Get the card from the database
	card, err := models.GetCard(db.(*gorm.DB), cardId)
	if err != nil {
		ctx.JSON(404, ErrInvalidCard)
		return
	}

	// Get the game engine from the context
	gameEngine, ok := ctx.Get(engine.EngineKey)
	if !ok {
		ctx.JSON(500, ErrInternalError)
		return
	}

	// Check if the game is finished
	if card.LastGame > gameEngine.(*engine.Engine).GetGameNumber() {
		log.Infof("Card last game %d and engine game %d", card.LastGame, gameEngine.(*engine.Engine).GetGameNumber())

		ctx.JSON(404, ErrUnfinishedGames)
		return
	}

	// Check the card
	amount := card.CheckCard(db.(*gorm.DB))

	// Return the ammount
	ctx.JSON(200, CheckCardResponse{Amount: amount})
}

type CheckCardResponse struct {
	Amount uint64 `json:"amount"`
}
