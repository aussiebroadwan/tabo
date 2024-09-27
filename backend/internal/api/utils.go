package api

type APIError struct {
	Message string `json:"message"`
}

var (
	ErrInvalidPicks    = APIError{Message: "Invalid picks"}
	ErrUnfinishedGames = APIError{Message: "Games haven't finished"}
	ErrInvalidCard     = APIError{Message: "Invalid Card ID"}
	ErrInternalError   = APIError{Message: "Internal Error"}
)
