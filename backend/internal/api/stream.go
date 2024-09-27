package api

import (
	"keno/internal/engine"
	"keno/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Stream Games Live
// @Summary Stream the current game
// @Description When a game is calculated and started, this endpoint will stream the game to the client. This will include all the picks which the client will have to display over 1.5 minutes for the proper effect.
// @Tags games
// @Produce json
// @Param Connection header string true "Connection: Upgrade"
// @Param Upgrade header string true "Upgrade: websocket"
// @Param Sec-Websocket-Version header string true "Sec-Websocket-Version: 13"
// @Success 200 {object} models.Message
// @Failure 500 {object} APIError
// @Router /api/v1/ws [get]
func GameStreamer(ctx *gin.Context) {
	// Get the game engine from the context
	gameEngine, ok := ctx.Get(engine.EngineKey)
	if !ok {
		ctx.JSON(500, ErrInternalError)
		return
	}

	// Upgrade to a websocket
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.WithError(err).Error("Error upgrading to websocket")
		return
	}
	defer ws.Close()

	// Add the listener
	listener := make(chan models.Message, 10)
	gameEngine.(*engine.Engine).AddListener(listener)
	defer gameEngine.(*engine.Engine).RemoveListener(listener)

	// While connection is open
	for {
		select {
		case message := <-listener:
			// Send the Game Data
			err := ws.WriteJSON(message)
			if err != nil {
				log.WithError(err).Error("Error writing game to websocket")
				return
			}
		case <-time.After(5 * time.Second):
			// Send a ping to keep the connection alive
			ws.WriteMessage(websocket.PingMessage, []byte{})
		case <-ctx.Done():
			return
		}
	}
}
