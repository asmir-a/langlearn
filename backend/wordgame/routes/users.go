package routes

import (
	"net/http"

	"github.com/asmir-a/gorestrouter"
	"github.com/asmir-a/langlearn/backend/auth/routes"
	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/logic"
)

func handlerBuilderStats(params map[string]string) http.Handler {
	username := params["username"]
	handlerWordgame := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
		wordCounts, httpErr := logic.GetWordCountsFor(username)
		if httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(wordCounts)
		return nil
	}
	return httperrors.HandlerWithHttpError(handlerWordgame)
}

func handlerBuilderWords(params map[string]string) http.Handler {
	username := params["username"]
	handlerWords := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
		wordsLists, httpErr := logic.GetWordsFor(username)
		if httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(wordsLists)
		return nil
	}
	return httperrors.HandlerWithHttpError(handlerWords)
}

func NewKnowsRouter() *gorestrouter.Router {
	//prolly it is fine if this router handles the authorization logic
	knowsRouter := &gorestrouter.Router{}

	authorizedHandlerBuilderStats := routes.CheckIfAuthorizedForUsername(handlerBuilderStats) //need to rethink about where all the authorizing logic and all the authentication logic should be set up; they can be separated from the application logic
	knowsRouter.Handle("/[username]/word-counts", authorizedHandlerBuilderStats)              //need to protect the router here

	authorizedHandlerBuilderWords := routes.CheckIfAuthorizedForUsername(handlerBuilderWords)
	knowsRouter.Handle("/[username]/words", authorizedHandlerBuilderWords)

	return knowsRouter
}
