package auth

import (
	"errors"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
	"github.com/asmir-a/langlearn/backend/auth/logic/passwords"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

func checkIfUserExistsForLogin(username string, password string) *httperrors.HttpError {
	userExists, httpErr := dbwrappers.CheckIfUserExists(username)
	if httpErr != nil {
		return httperrors.NewHttp500Error(httpErr)
	}
	if !userExists {
		return httperrors.NewHttpError(
			errors.New("user with that username does not exist"),
			http.StatusUnauthorized,
			"user with provided username does not exist",
		)
	}
	return nil
}

func isPasswordCorrect(username string, password string) *httperrors.HttpError {
	salt, httpErr := dbwrappers.GetUserPasswordSalt(username)
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	potentialPasswordHash := passwords.Hash(password, salt)
	validPasswordHash, httpErr := dbwrappers.GetUserPasswordHash(username)
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	if potentialPasswordHash != validPasswordHash {
		return httperrors.NewHttpError(
			errors.New("wrong password"),
			http.StatusUnauthorized,
			"wrong password",
		)
	}
	return nil
}

func createNewSessionFor(username string) (string, *httperrors.HttpError) {
	sessionExists, httpErr := dbwrappers.CheckIfSessionExistsFor(username)
	if httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	}
	if !sessionExists {
		sessionKey, httpErr := dbwrappers.CreateSessionFor(username)
		if httpErr != nil {
			return "", httperrors.WrapError(httpErr)
		}
		return sessionKey, nil
	}
	sessionKey, httpErr := dbwrappers.ReplaceSessionFor(username)
	if httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	}
	return sessionKey, nil
}

func Login(username string, password string) (string, *httperrors.HttpError) {
	if httpErr := validateCredentials(username, password); httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	}

	if httpErr := checkIfUserExistsForLogin(username, password); httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	}

	if httpErr := isPasswordCorrect(username, password); httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	}

	if sessionKey, httpErr := createNewSessionFor(username); httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	} else {
		return sessionKey, nil
	}
	//todo: the session should be checked and deleted in a single database transaction
}
