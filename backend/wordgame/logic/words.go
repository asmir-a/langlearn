package logic

import (
	"encoding/json"

	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/dbwrappers"
)

type WordsLists struct {
	Learned  []string `json:"learned"`
	Learning []string `json:"learning"`
}

func getWordsLearnedFor(username string) ([]string, *httperrors.HttpError) {
	wordsLearned, httpErr := dbwrappers.SelectWordsLearnedFor(username)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}
	return wordsLearned, nil
}

func getWordsLearningFor(username string) ([]string, *httperrors.HttpError) {
	wordsLearning, httpErr := dbwrappers.SelectWordsLearningFor(username)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}
	return wordsLearning, nil
}

func GetWordsFor(username string) ([]byte, *httperrors.HttpError) {
	wordsLearned, httpErr := getWordsLearnedFor(username)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}
	wordsLearning, httpErr := getWordsLearningFor(username)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}

	var words WordsLists
	words.Learned = wordsLearned
	words.Learning = wordsLearning

	wordsJson, err := json.Marshal(&words)
	if err != nil {
		return nil, httperrors.NewHttp500Error(err)
	}

	return wordsJson, nil
}
