package routes

import (
	"net/http"

	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/logic"
)

func handleGameEntriesRandom(w http.ResponseWriter, _ *http.Request) *httperrors.HttpError {
	gameEntryJson, httpErr := logic.GetGameEntryJson()
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(gameEntryJson))
	return nil
}

func handleGameEntriesSubmit(w http.ResponseWriter, req *http.Request) {
}
