package auth

import (
	"errors"
	"net/http"

	"github.com/asmir-a/langlearn/backend/httperrors"
)

//todo: this must reside somewhere else
//maybe logic/validation
//or logic/utilities

func validateUsername(username string) bool {
	//may be better to have a separate file for validation
	if username == "" {
		return false
	}
	//todo: other checks to prevent sql injection eg
	return true
}

func validatePassword(password string) bool {
	if password == "" {
		return false
	}
	return true
}

func validateCredentials(username string, password string) *httperrors.HttpError {
	if !validateUsername(username) || !validatePassword(password) {
		return httperrors.NewHttpError( //any instances of new errors should prolly be inside of functions
			errors.New("invalid login or password"),
			http.StatusUnauthorized,
			"invalid username or password",
		)
	}
	return nil
}
