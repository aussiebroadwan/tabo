package models

import (
	"keno/internal/utils"

	"gorm.io/gorm"
)

type Game struct {
	ID    uint64  `json:"id" gorm:"primary_key"`
	Picks []uint8 `json:"picks"`
}

// CheckGame is a method that checks the game for matches against the selection
// and returns the number of matches. This doesn't handle the payout, that is
// done in the payout package.
func (g Game) CheckGame(selection []uint8) uint8 {
	// Rolling Amount
	matches := uint8(0)

	for _, pick := range g.Picks {
		if utils.Contains(selection, pick) {
			matches++
		}
	}

	return matches
}

func GetGame(db *gorm.DB, id uint64) (*Game, error) {
	var game Game
	err := db.First(&game, id).Error
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func GetLastGame(db *gorm.DB) (*Game, error) {
	var game Game
	err := db.Last(&game).Error
	if err != nil {
		return nil, err
	}

	return &game, nil
}

// CommitNewGame is a method that commits a new game to the database. You don't
// need to have any picks to commit a new game, but you will need to commit
// all 20 picks before you can check the game properly.
func CommitNewGame(db *gorm.DB, game *Game) error {
	tx := db.Create(game)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// CommitGamePick is a method that commits a new pick to the database. You will
// need to commit all 20 picks before you can check the game properly.
func CommitGamePick(db *gorm.DB, game *Game, pick uint8) error {
	tx := db.Model(game).Update("picks", append(game.Picks, pick))
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// Pulled directly from https://www.keno.com.au/keno-pdfs/NSW_Game%20Guide.pdf
//
// It contains all the payouts for the game based on then number of picks the
// client has made and the number of matches they have.
//
//	eg. If the client has made 10 picks and has 5 matches, they will win $2
//			matchMatrix[10][5] = 2
var matchMatrix map[uint8]map[uint8]uint64 = map[uint8]map[uint8]uint64{
	// Small Picks
	1: {1: 3},
	2: {2: 12},
	3: {3: 44, 2: 1},
	4: {4: 120, 3: 4, 2: 1},
	5: {5: 640, 4: 14, 3: 2},
	6: {6: 1_800, 5: 80, 4: 5, 3: 1},

	// Jackpots
	7:  {7: 5_000, 6: 125, 5: 12, 4: 3, 3: 1},
	8:  {8: 25_000, 7: 675, 6: 60, 5: 7, 4: 2},
	9:  {9: 100_000, 8: 2_500, 7: 210, 6: 20, 5: 5, 4: 1},
	10: {10: 1_000_000, 9: 10_000, 8: 580, 7: 50, 6: 6, 5: 2, 4: 1},

	// Big Picks
	15: {15: 250_000, 14: 100_000, 13: 50_000, 12: 12_000, 11: 2_000, 10: 250, 9: 50, 8: 20, 7: 4, 6: 2, 5: 1},
	20: {20: 250_000, 19: 100_000, 18: 50_000, 17: 25_000, 16: 15_000, 15: 10_000, 14: 5_000, 13: 1_200, 12: 450, 11: 100, 10: 20, 9: 7, 8: 2, 2: 2, 1: 10, 0: 100},
	40: {20: 250_000, 19: 25_000, 18: 2_200, 17: 200, 16: 35, 15: 7, 14: 2, 13: 1, 7: 1, 6: 2, 5: 7, 4: 35, 3: 200, 2: 2_200, 1: 25_000, 0: 250_000},
}

func winMatrix(numsPlayed, numsMatched uint8) uint64 {
	section, ok := matchMatrix[numsPlayed]
	if !ok {
		return 0
	}

	amount, ok := section[numsMatched]
	if !ok {
		return 0
	}

	return amount
}
