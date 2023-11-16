package models

import (
	"sort"
	"time"

	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

type Card struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time

	Selection []uint8 `json:"selection"`
	StartGame uint64  `json:"start_game_num"`
	LastGame  uint64  `json:"last_game_num"`
	PerGame   uint64  `json:"per_game"`

	User string `json:"user"`
}

func (c Card) CheckCard(db *gorm.DB) uint64 {
	// Rolling Amount
	amount := uint64(0)

	for gameNum := c.StartGame; gameNum < c.LastGame; gameNum++ {
		// Get the game
		game, err := GetGame(db, gameNum)
		if err != nil {
			log.WithError(err).Error("Error getting game")
			continue
		}

		// Check the game
		matches := game.CheckGame(c.Selection)
		amount += winMatrix(uint8(len(c.Selection)), matches) * c.PerGame

	}

	return amount
}

func GetCard(db *gorm.DB, id uint64) (*Card, error) {
	var card Card
	err := db.First(&card, id).Error
	if err != nil {
		return nil, err
	}

	return &card, nil
}

func SubmitCard(db *gorm.DB, selection []uint8, startGame uint64, numOfGames uint8, pricePerGame uint64, user string) (*Card, error) {
	// Sort selection
	sort.Slice(selection, func(i, j int) bool { return selection[i] < selection[j] })

	// Setup the Card
	newCard := &Card{
		CreatedAt: time.Now(),
		Selection: selection,
		StartGame: startGame,
		LastGame:  startGame + uint64(numOfGames),
		PerGame:   uint64(pricePerGame),
		User:      user,
	}

	// Create the card in the database
	tx := db.Create(newCard)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return newCard, nil
}
