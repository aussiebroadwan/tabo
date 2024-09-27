package models

// Message is a message that is sent to the client over the websocket, it
// contains the message type and the body of the message and is used to
// communicate current game state with the client.
type Message struct {
	// The message type
	Type string        `json:"type"`
	Body StreamMessage `json:"body"`
}

type StreamMessage interface {
	GetType() string
}

func GenerateMessage(msg StreamMessage) Message {
	return Message{
		Type: msg.GetType(),
		Body: msg,
	}
}

// NewGame is a message that is sent to the client when a new game is started,
// it contains the game id, the next game time, the current game start time,
// and the current game end time. It is sent once per game and indicates to the
// client that it should reset its state.
type NewGameMsg struct {
	GameId               uint64 `json:"gameId"`
	NextGameTime         int64  `json:"nextGameTime"`
	CurrentGameStartTime int64  `json:"currentGameStartTime"`
	CurrentGameEndTime   int64  `json:"currentGameEndTime"`
}

func (n NewGameMsg) GetType() string {
	return "NEW"
}

// NewPick is a message that is sent to the client when a new pick is made in
// the game, it contains the pick number. There will be 20 of these messages
// sent per game.
type NewPickMsg struct {
	Pick int `json:"pick"`
}

func (n NewPickMsg) GetType() string {
	return "PIC"
}

// CurrentGame is a message that is sent to the client when it first connects
// to the websocket, it contains the game id, the game times, and the picks
// that have been made so far in the game. It is sent once per connection and
// indicates to the client what the current game state is.
type CurrentGameMsg struct {
	GameId               uint64 `json:"gameId"`
	NextGameTime         int64  `json:"nextGameTime"`
	CurrentGameStartTime int64  `json:"currentGameStartTime"`
	CurrentGameEndTime   int64  `json:"currentGameEndTime"`
	Picks                []int  `json:"picks"`
}

func (c CurrentGameMsg) GetType() string {
	return "CUR"
}
