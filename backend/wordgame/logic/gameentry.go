package logic

import (
	"encoding/json"
	"errors"

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

func HandleAnswer(submission WordGameSubmission) *httperrors.HttpError {
	if submission.IsAnswerCorrect {
		httpErr := handleCorrectAnswer(submission.Username, submission.Word)
		if httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		return nil
	}

	httpErr := handleIncorrectAnswer(submission.Username, submission.Word)
	if httpErr != nil {
		return httperrors.WrapError(httpErr) //todo: let WrapError handle nil errors so that you could just return httperrors.WrapError(httpErr)
	}
	return nil
}

func handleCorrectAnswer(username string, word string) *httperrors.HttpError {
	rowExists, httpErr := dbwrappers.DoesRowExists(username, word)
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	if !rowExists {
		httpErr = dbwrappers.CreateNewRowInKnows(username, word)
		if httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		return nil
	}

	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	httpErr = dbwrappers.IncrementCurrentCount(username, word)
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	return nil
}

func handleIncorrectAnswer(username string, word string) *httperrors.HttpError {
	rowExists, httpErr := dbwrappers.DoesRowExists(username, word)
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}

	if !rowExists {
		return nil
	}

	currentCount, httpErr := dbwrappers.GetCurrentCount(username, word)
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}

	if currentCount <= 0 { //should be done by db contraints instead
		return httperrors.NewHttp500Error(errors.New("current count cannot be equal to or less than 0"))
	}
	if currentCount == 1 {
		httpErr := dbwrappers.DeleteRowInKnows(username, word)
		if httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		return nil
	}
	httpErr = dbwrappers.DecrementCurrentCount(username, word)
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	return nil
}
