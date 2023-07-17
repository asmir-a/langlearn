package routes

import (
	"encoding/json"
	"net/http"

	"github.com/asmir-a/gorestrouter"
	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/logic"
)

func NewGameEntriesRouter() *gorestrouter.Router {
	router := gorestrouter.Router{}

	//the authentication and authorization logic should be handled somewhere here probably
	router.Handle("/users/[username]/random", handlerBuilderGameEntriesRandom)
	router.Handle("/users/[username]/submit", handlerBuilderGameEntriesSubmit)

	return &router
}

func handlerBuilderGameEntriesRandom(params map[string]string) http.Handler {
	username := params["username"]
	handlerGameEntries := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
		gameEntryJson, httpErr := logic.GetGameEntry(username)
		if httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(gameEntryJson))
		return nil
	}
	return httperrors.HandlerWithHttpError(handlerGameEntries)
}

func extractSubmissionFromRequest(req *http.Request) (logic.WordGameSubmission, *httperrors.HttpError) { //this should be refactored: the submission could be extracted from the body intead and this function should be in the submission.go file; then, the struct would not need to be exported
	decoder := json.NewDecoder(req.Body)
	var submission logic.WordGameSubmission
	if err := decoder.Decode(&submission); err != nil {
		return logic.WordGameSubmission{}, httperrors.NewHttp500Error(err)
	}
	return submission, nil
}

func handlerBuilderGameEntriesSubmit(params map[string]string) http.Handler {
	username := params["username"]
	handlerGameEntriesSubmit := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
		submission, httpErr := extractSubmissionFromRequest(req)
		if httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		if httpErr = logic.HandleAnswer(username, submission); httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		return nil
	}
	return httperrors.HandlerWithHttpError(handlerGameEntriesSubmit)
}
