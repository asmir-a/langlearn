package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
	"github.com/asmir-a/langlearn/backend/auth/logic/passwords"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

// todo: need to refactor this function; it is too long
func Login(username string, password string) (string, *httperrors.HttpError) {
	if !validateUsername(username) || !validatePassword(password) {
		return "", httperrors.NewHttpError(
			errors.New("invalid username or password"),
			http.StatusUnauthorized,
			"invalid username or password format",
			currentFilePath+fmt.Sprintf("Login:!validate with (username, password) = (%s,%s)", username, password),
		)
	}

	userExists, err := dbwrappers.CheckIfUserExists(username)
	if err == nil && !userExists {
		return "", httperrors.NewHttpError(
			errors.New("user with that username does not exist"),
			http.StatusUnauthorized,
			"user with provided username does not exist",
			currentFilePath+fmt.Sprintf("Login:checkIfUserExists with username = %s", username),
		)
	} else if err != nil {
		return "", httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			currentFilePath+fmt.Sprintf("Login:checkIfUserExists with username = %s", username),
		)
	}

	salt, err := dbwrappers.GetUserPasswordSalt(username)
	if err != nil {
		return "", httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			currentFilePath+fmt.Sprintf("Login:getUserPasswordSalt with username = %s", username), //we cannot append to debug because this is the first error source
		)
	}

	potentialPasswordHash := passwords.Hash(password, salt)
	validPasswordHash, err := dbwrappers.GetUserPasswordHash(username)
	if err != nil {
		return "", httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			currentFilePath+fmt.Sprintf("Login:getUserPasswordSalt with username = %s", username), //todo: status code and message can be encapsulated in NewInternalServerError() function and the current path can be encapsulated in a multilevel closure: one level for current file, one level for current function, and finally we just pass the string indicating where in the function the error happened
		)
	}

	if potentialPasswordHash != validPasswordHash {
		//todo: may be need to delete the session or not
		return "", httperrors.NewHttpError(
			errors.New("wrong password"),
			http.StatusUnauthorized,
			"wrong password",
			currentFilePath+fmt.Sprintf("Login:potentialPasswordHash!=validPasswordHash"),
		)
	}

	sessionExists, err := dbwrappers.CheckIfSessionExistsFor(username) //might be unncessary; we can set up the constraint in the database
	if err != nil {
		return "", httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			currentFilePath+"Login:checkIfSessionExists",
		)
	}
	if !sessionExists {
		sessionKey, err := dbwrappers.CreateSessionFor(username)
		if err != nil {
			return "", httperrors.NewHttpError(
				err,
				http.StatusInternalServerError,
				"something went wrong",
				currentFilePath+"Login:createSessionFor",
			)
		}
		return sessionKey, nil
	}

	sessionKey, err := dbwrappers.ReplaceSessionFor(username) //we do not care if the old session is valid or not; since we got the right login and password, we need to create a new valid session
	if err != nil {
		return "", httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			currentFilePath+fmt.Sprintf("Login:replaceSessionFor with username=%s", username),
		)
	}
	return sessionKey, nil
	//todo: the session should be checked and deleted in a single database transaction
}
