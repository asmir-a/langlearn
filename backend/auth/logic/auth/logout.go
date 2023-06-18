package auth

import (
	"errors"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

func checkIfSessionIsValid(sessionKey string) *httperrors.HttpError {
	if sessionValid, httpErr := dbwrappers.CheckIfSessionIsValid(sessionKey); httpErr != nil {
		return httperrors.WrapError(httpErr)
	} else if !sessionValid {
		return httperrors.NewHttpError(
			errors.New("session invalid"),
			http.StatusUnauthorized,
			"please log in",
		)
	}
	return nil
}

func Logout(currentSessionKey string) *httperrors.HttpError {
	if httpErr := checkIfSessionIsValid(currentSessionKey); httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	if httpErr := dbwrappers.DeleteSession(currentSessionKey); httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	return nil
}
