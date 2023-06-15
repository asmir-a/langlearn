package logic

import (
	"encoding/json"

	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/dbwrappers"
	"github.com/asmir-a/langlearn/backend/wordgame/imagesearch"
)

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

func GetGameEntryJson() (string, *httperrors.HttpError) {
	gameEntry, httpErr := getGameEntry()
	if httpErr != nil {
		return "", httpErr
	}

	gameEntryBytes, err := json.Marshal(gameEntry)
	if err != nil {
		newHttpErr := httperrors.NewHttp500Error(err)
		return "", newHttpErr
	}

	return string(gameEntryBytes), nil
}
