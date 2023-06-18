package auth

import (
	"errors"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
	"github.com/asmir-a/langlearn/backend/auth/logic/passwords"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

func checkIfUserExistsForSignup(username string) *httperrors.HttpError {
	usernameExists, httpErr := dbwrappers.CheckIfUserExists(username)

	if httpErr == nil && usernameExists {
		return httperrors.NewHttpError(
			errors.New("username already exists"),
			http.StatusConflict,
			"username already exists",
		)
	}

	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}

	return nil
}

func Signup(username string, password string) (string, *httperrors.HttpError) {
	if httpErr := validateCredentials(username, password); httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	}

	if httpErr := checkIfUserExistsForSignup(username); httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	}

	salt, httpErr := passwords.Salt(username)
	if httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	}
	hash := passwords.Hash(password, salt)

	if httpErr = dbwrappers.InsertUser(username, hash, salt); httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	}

	sessionKey, httpErr := dbwrappers.CreateSessionFor(username)
	if httpErr != nil {
		return "", httperrors.WrapError(httpErr) //todo: use reflection to say that httperror cannot be used in wraperror
	}

	return sessionKey, nil
}
