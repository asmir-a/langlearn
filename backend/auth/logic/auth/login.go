package auth

import (
	"errors"
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
		)
	}

	userExists, err := dbwrappers.CheckIfUserExists(username)
	if err == nil && !userExists {
		return "", httperrors.NewHttpError(
			errors.New("user with that username does not exist"),
			http.StatusUnauthorized,
			"user with provided username does not exist",
		)
	} else if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}

	salt, err := dbwrappers.GetUserPasswordSalt(username)
	if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}

	potentialPasswordHash := passwords.Hash(password, salt)
	validPasswordHash, err := dbwrappers.GetUserPasswordHash(username)
	if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}

	if potentialPasswordHash != validPasswordHash {
		//todo: may be need to delete the session or not
		return "", httperrors.NewHttpError(
			errors.New("wrong password"),
			http.StatusUnauthorized,
			"wrong password",
		)
	}

	sessionExists, err := dbwrappers.CheckIfSessionExistsFor(username) //might be unncessary; we can set up the constraint in the database
	if err != nil {
		return "", httperrors.NewHttp500Error(err)

	}
	if !sessionExists {
		sessionKey, err := dbwrappers.CreateSessionFor(username)
		if err != nil {
			return "", httperrors.NewHttp500Error(err)
		}
		return sessionKey, nil
	}

	sessionKey, err := dbwrappers.ReplaceSessionFor(username) //we do not care if the old session is valid or not; since we got the right login and password, we need to create a new valid session
	if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}
	return sessionKey, nil
	//todo: the session should be checked and deleted in a single database transaction
}
