package routes

import (
	"net/http"

	"github.com/asmir-a/gorestrouter"
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

	knowsRouter.Handle("/[username]/word-counts", handlerBuilderStats) //need to protect the router here
	knowsRouter.Handle("/[username]/words", handlerBuilderWords)

	return knowsRouter
}
