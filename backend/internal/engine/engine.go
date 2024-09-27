package engine

import (
	"fmt"
	"keno/internal/models"
	"math/rand"
	"sync"
	"time"

	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

const (
	PlayTime = 90 * time.Second
	WaitTime = 90 * time.Second

	NumberPicks    = 20
	NumberRangeMin = 1
	NumberRangeMax = 80

	EngineKey = "engine"
)

type Engine struct {
	gameNumber      uint64
	nextGameTime    time.Time
	curGamStartTime time.Time
	curGame         models.Game
	mu              sync.RWMutex
	db              *gorm.DB

	listeners []chan models.Message
}

func SetupEngine(db *gorm.DB) *Engine {
	var activeGameNum uint64 = 1

	// Get Last Game if it exists
	game, err := models.GetLastGame(db)
	if err == nil {
		activeGameNum = game.ID + 1
	}

	return &Engine{
		gameNumber:      activeGameNum,
		nextGameTime:    time.Now(),
		curGamStartTime: time.Now(),
		curGame:         models.Game{},
		db:              db,
		mu:              sync.RWMutex{},
		listeners:       make([]chan models.Message, 0),
	}
}

// ==================
// 		Getters
// ==================

// GetGameNumber is a method that returns the current game number.
// This method is useful when you want to fetch the current game number
// from the engine instance in a thread-safe manner.
func (engine *Engine) GetGameNumber() uint64 {
	engine.mu.RLock()
	defer engine.mu.RUnlock()

	return engine.gameNumber
}

// GetNextGame is a method that returns the time of the next game.
// This method is useful when you want to know when the next game will start
// from the engine instance in a thread-safe manner.
func (engine *Engine) GetNextGame() time.Time {
	engine.mu.RLock()
	defer engine.mu.RUnlock()

	return engine.nextGameTime
}

// GetCurGameStart is a method that returns the start time of the current game.
// This method is useful when you want to know when the current game started
// from the engine instance in a thread-safe manner.
func (engine *Engine) GetCurGameStart() time.Time {
	engine.mu.RLock()
	defer engine.mu.RUnlock()

	return engine.curGamStartTime
}

// ==================
// Notification Logic
// ==================

// AddListener adds a new listener to the engine. Listeners are clients who
// are connected via websockets and need to get game state updates. It also
// sends the current game message to the newly added listener. This is
// particularly useful when a client connects and needs to know the state of
// the current game immediately.
func (engine *Engine) AddListener(listener chan models.Message) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	engine.listeners = append(engine.listeners, listener)

	// Generate a Current Game Message
	curGame := models.CurrentGameMsg{
		GameId:               engine.curGame.ID,
		NextGameTime:         engine.nextGameTime.UnixMilli(),
		CurrentGameStartTime: engine.curGamStartTime.UnixMilli(),
		CurrentGameEndTime:   engine.curGamStartTime.Add(PlayTime).UnixMilli(),
		Picks:                make([]int, 0),
	}

	// Converting picks to ints
	for _, pick := range engine.curGame.Picks {
		curGame.Picks = append(curGame.Picks, int(pick))
	}

	// Send current game
	select {
	case listener <- models.GenerateMessage(curGame):
	default:
		log.WithField("src", "engine.AddListener").Error("Listener channel full")
	}
}

// RemoveListener is a method that removes an existing listener from the engine.
// Disconnected clients should be removed from the engine so that the engine
// doesn't try to send messages to a closed channel.
func (engine *Engine) RemoveListener(listener chan models.Message) {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	for i, l := range engine.listeners {
		if l == listener {
			engine.listeners = append(engine.listeners[:i], engine.listeners[i+1:]...)
			return
		}
	}
}

// This method is useful when you want to notify all the listeners about a new
// game message. This method is called whenever the game state changes.
func (engine *Engine) NotifyListeners(game models.Message) {
	engine.mu.RLock()
	defer engine.mu.RUnlock()

	for _, listener := range engine.listeners {
		select {
		case listener <- game:
		default:
			log.WithField("src", "engine.NotifyListeners").Error("Listener channel full")
		}
	}
}

// ==================
//     Game Logic
// ==================

func (engine *Engine) StartLoop() {
	for {
		game := engine.initialiseGame()

		// Generate Random Number Generator
		picks := map[int]bool{}

		// Generate Picks
		for i := 0; i < NumberPicks; i++ {
			engine.generatePick(game, &picks, i)
			time.Sleep(engine.calculatePickSleepDuration(i))
		}

		// Increment Game Number
		engine.mu.Lock()
		engine.gameNumber++
		engine.mu.Unlock()

		// Sleep till next game
		time.Sleep(time.Until(engine.nextGameTime))

		log.WithFields(log.Fields{
			"src":   "engine.StartLoop",
			"game":  game.ID,
			"picks": fmt.Sprintf("%+v", game.Picks),
		}).Info("Game Complete")
	}
}

func (engine *Engine) initialiseGame() *models.Game {
	// Set the game times for the new game
	engine.mu.Lock()
	engine.curGamStartTime = time.Now()
	engine.nextGameTime = engine.curGamStartTime.Add(PlayTime + WaitTime)
	engine.mu.Unlock()

	// Start new Game

	game := &models.Game{
		ID:    engine.gameNumber,
		Picks: []uint8{},
	}

	// Commit the Game to storage and notify listeners to clear state
	// and get ready for the next game.
	models.CommitNewGame(engine.db, game)
	engine.NotifyListeners(models.GenerateMessage(models.NewGameMsg{
		GameId:               game.ID,
		NextGameTime:         engine.nextGameTime.UnixMilli(),
		CurrentGameStartTime: engine.curGamStartTime.UnixMilli(),
		CurrentGameEndTime:   engine.curGamStartTime.Add(PlayTime).UnixMilli(),
	}))
	engine.mu.Lock()
	engine.curGame = *game
	engine.mu.Unlock()

	return game
}

func (engine *Engine) generatePick(game *models.Game, picks *map[int]bool, i int) {
	seed := uint64(time.Now().UnixNano()) + rand.Uint64()
	r := rand.New(rand.NewSource(int64(seed)))

	// Pick a number that hasn't been picked
	pick := 0
	for (*picks)[pick] || pick == 0 {
		pick = r.Int()%(NumberRangeMax-NumberRangeMin+1) + NumberRangeMin
	}

	// Add pick to map
	(*picks)[pick] = true
	game.Picks = append(game.Picks, uint8(pick))

	// Commit the pick to storage and notify listeners
	models.CommitGamePick(engine.db, game, uint8(pick))
	engine.NotifyListeners(models.GenerateMessage(models.NewPickMsg{
		Pick: pick,
	}))
	engine.mu.Lock()
	engine.curGame = *game
	engine.mu.Unlock()
}

func (engine *Engine) calculatePickSleepDuration(i int) time.Duration {
	gameEndTime := engine.curGamStartTime.Add(PlayTime)
	return time.Until(gameEndTime) / time.Duration(NumberPicks-i)
}
