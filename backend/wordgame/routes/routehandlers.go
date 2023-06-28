package routes

import (
	"encoding/json"
	"net/http"

	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/logic"
)

func handleGameEntriesRandom(w http.ResponseWriter, _ *http.Request) *httperrors.HttpError {
	gameEntryJson, httpErr := logic.GetGameEntry()
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(gameEntryJson))
	return nil
}

func extractSubmissionFromRequest(req *http.Request) (logic.WordGameSubmission, *httperrors.HttpError) {
	decoder := json.NewDecoder(req.Body)
	var submission logic.WordGameSubmission
	if err := decoder.Decode(&submission); err != nil {
		return logic.WordGameSubmission{}, httperrors.NewHttp500Error(err)
	}
	return submission, nil
}

func handleGameEntriesSubmit(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
	submission, httpErr := extractSubmissionFromRequest(req)
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	if httpErr = logic.HandleAnswer(submission); httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	return nil
}
