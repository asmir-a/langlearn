package routes

import (
	"net/http"

	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/logic"
)

const fileLevelDebugInfo = "wordgame routes"

func handleGameEntriesRandom(w http.ResponseWriter, _ *http.Request) *httperrors.HttpError {
	const funcLevelDebugInfo = "HandleGameEntry "
	gameEntryJson, httpErr := logic.GetGameEntryJson()
	if httpErr != nil {
		return httperrors.WrapError(httpErr, fileLevelDebugInfo+funcLevelDebugInfo)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(gameEntryJson))
	return nil
}

func handleGameEntriesSubmit(w http.ResponseWriter, req *http.Request) {
}
