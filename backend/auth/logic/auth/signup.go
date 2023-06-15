package auth

import (
	"errors"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
	"github.com/asmir-a/langlearn/backend/auth/logic/passwords"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

// todo: need to refactor; split to multiple functions
func Signup(username string, password string) (string, *httperrors.HttpError) {
	if !validateUsername(username) || !validatePassword(password) {
		return "", httperrors.NewHttpError(
			errors.New("invalid login or password"),
			http.StatusUnauthorized,
			"invalid username or password",
		)
	}

	usernameExists, err := dbwrappers.CheckIfUserExists(username)
	if err == nil && usernameExists {
		return "", httperrors.NewHttpError(
			errors.New("username already exists"),
			http.StatusConflict,
			"username already exists",
		)
	} else if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}

	salt, err := passwords.Salt(username)
	if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}

	hash := passwords.Hash(password, salt)
	if err = dbwrappers.InsertUser(username, hash, salt); err != nil {
		return "", httperrors.NewHttp500Error(err)
	}

	sessionKey, err := dbwrappers.CreateSessionFor(username)
	if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}

	return sessionKey, nil
}
