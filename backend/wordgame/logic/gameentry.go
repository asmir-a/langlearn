package logic

import (
	"encoding/json"

	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/dbwrappers"
	"github.com/asmir-a/langlearn/backend/wordgame/imagesearch"
)

type WordGameEntry struct {
	CorrectWord         string   `json:"correctWord"`
	CorrectWordImageUrl string   `json:"correctWordImageUrl"`
	IncorrectWords      []string `json:"incorrectWords"`
}

func getGameEntry() (WordGameEntry, *httperrors.HttpError) {
	words, httpErr := dbwrappers.GetRandomKoreanWords()
	if httpErr != nil {
		return WordGameEntry{}, httperrors.WrapError(httpErr)
	}

	correctWordIndex := 0

	correctWord := words[correctWordIndex].Word
	correctWordDef := words[correctWordIndex].Defs[0]

	correctWordImageUrl, httpErr := imagesearch.FetchImageUrlFor(correctWordDef)
	if httpErr != nil {
		return WordGameEntry{}, httperrors.WrapError(httpErr)
	}

	incorrectWordsWithDefs := words[1:]
	incorrectWords := dbwrappers.ExtractWordsFromWordsWithDefs(incorrectWordsWithDefs) //todo: this does not belong to dbwrappers; also the type prolly should be shared among these packages

	return WordGameEntry{
		CorrectWord:         correctWord,
		CorrectWordImageUrl: correctWordImageUrl,
		IncorrectWords:      incorrectWords,
	}, nil
}

func GetGameEntry() ([]byte, *httperrors.HttpError) {
	gameEntry, httpErr := getGameEntry()
	if httpErr != nil {
		return nil, httpErr
	}

	gameEntryBytes, err := json.Marshal(gameEntry)
	if err != nil {
		newHttpErr := httperrors.NewHttp500Error(err)
		return nil, newHttpErr
	}

	return gameEntryBytes, nil
}
