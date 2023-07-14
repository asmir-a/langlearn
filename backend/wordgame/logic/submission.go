package logic

import (
	"errors"

	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/dbwrappers"
)

type WordGameSubmission struct {
	IsAnswerCorrect bool   `json:"isAnswerCorrect"`
	Word            string `json:"word"`
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

func HandleAnswer(username string, submission WordGameSubmission) *httperrors.HttpError {
	if submission.IsAnswerCorrect {
		httpErr := handleCorrectAnswer(username, submission.Word)
		if httpErr != nil {
			return httperrors.WrapError(httpErr)
		}
		return nil
	}

	httpErr := handleIncorrectAnswer(username, submission.Word)
	if httpErr != nil {
		return httperrors.WrapError(httpErr) //todo: let WrapError handle nil errors so that you could just return httperrors.WrapError(httpErr)
	}
	return nil
}
