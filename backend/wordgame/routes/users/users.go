package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/asmir-a/langlearn/backend/auth/routes"
	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/dbwrappers"
)

func getHead(currentPath string) string {
	currentPath = path.Clean("/"+currentPath) + "/"
	slashIndex := strings.Index(currentPath[1:], "/") + 1
	if slashIndex == 0 {
		return ""
	}
	return currentPath[1:slashIndex]
}

func getTail(currentPath string) string {
	currentPath = path.Clean("/" + currentPath)
	slashIndex := strings.Index(currentPath[1:], "/") + 1
	if slashIndex == 0 {
	}
	return currentPath[slashIndex:]
}

func getHeadAndTail(currentPath string) (string, string) {
	currentPath = path.Clean("/" + currentPath)
	if currentPath == "/" { //handles "/"" and "" and other weird things like "//////"
		return "", ""
	}
	currentPath = currentPath + "/"
	slashIndex := strings.Index(currentPath[1:], "/") + 1
	return currentPath[1:slashIndex], currentPath[slashIndex+1:]
}

func shiftPath(currentPath string) (string, string) {
	currentPath = path.Clean("/" + currentPath)
	slashIndex := strings.Index(currentPath[1:], "/") + 1
	if slashIndex == 0 { //strings.Index returns -1 if not found and we are adding 1 so 0
		return currentPath, "" //empty string signifies the end
	}
	return currentPath[1:slashIndex], currentPath[slashIndex+1:]
}

func preparePathForUsersRoute(req *http.Request) {
	expression := regexp.MustCompile(`/api/users/(?P<UsernameAndAfter>.*)`)
	currentPath := req.URL.Path
	matches := expression.FindStringSubmatch(currentPath)
	fmt.Println("matches[1]: ", matches[1])
	req.URL.Path = matches[1]
}

func usersRoutePreparer(currentHandler httperrors.HandlerWithHttpError) httperrors.HandlerWithHttpError {
	newFuncHandler := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
		preparePathForUsersRoute(req)
		if httpErr := currentHandler(w, req); httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		return nil
	}
	return httperrors.HandlerWithHttpError(newFuncHandler)
}

type WordCounts struct {
	LearnedCount  int `json:"learnedCount"`
	LearningCount int `json:"learningCount"`
}

func getWordCountsJsonFor(username string) ([]byte, *httperrors.HttpError) {
	learnedCount, httpErr := dbwrappers.GetWordsLearnedCount(username)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}
	learningCount, httpErr := dbwrappers.GetWordsLearningCount(username)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}
	wordCounts := WordCounts{
		LearnedCount:  learnedCount,
		LearningCount: learningCount,
	}
	byteBuf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(byteBuf)
	encoder.Encode(wordCounts)
	return byteBuf.Bytes(), nil
}

func wordGameHandlerWrapper(username string) httperrors.HandlerWithHttpError {
	wordGameHandler := func(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
		wordCountsJson, httpErr := getWordCountsJsonFor(username)
		fmt.Println(string(wordCountsJson))
		if httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(wordCountsJson)
		return nil
	}
	return httperrors.HandlerWithHttpError(wordGameHandler)
}

func usernameRouteHandler(w http.ResponseWriter, req *http.Request) *httperrors.HttpError {
	fmt.Println("usernameRouteHandler is called")
	username, tail := getHeadAndTail(req.URL.Path)
	fmt.Println("head: ", username)
	fmt.Println("tail: ", tail)
	if username == "" {
		return httperrors.NewHttpError(
			errors.New("trying to access unauthorized route"),
			http.StatusUnauthorized,
			"the route is not accessible",
		)
	}

	head, tail := getHeadAndTail(tail)
	fmt.Println("headnew: ", head)
	fmt.Println("tailnew: ", tail)

	switch head {
	case "wordgame":
		if httpErr := routes.CheckIfAuthorized(username, wordGameHandlerWrapper(username))(w, req); httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
	default:
		return httperrors.NewHttpError(
			errors.New("the resource is not found"),
			http.StatusNotFound,
			"the url is not accessible",
		)
	}

	return nil
}

func SetUpUsersWordsGameRoutes(mux *http.ServeMux) {
	mux.Handle("/api/users/", routes.CheckIfAuthenticated(
		usersRoutePreparer(
			httperrors.HandlerWithHttpError(
				usernameRouteHandler,
			),
		),
	)) //todo: handle /api/users path as well somehow
}
