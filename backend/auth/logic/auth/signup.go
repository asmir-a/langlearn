package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
	"github.com/asmir-a/langlearn/backend/auth/logic/passwords"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

func Signup(username string, password string) (string, *httperrors.HttpError) {
	if !validateUsername(username) || !validatePassword(password) {
		return "", httperrors.NewHttpError(
			errors.New("invalid login or password"),
			http.StatusUnauthorized,
			"invalid username or password",
			currentFilePath+fmt.Sprintf("Singup:!validate() with (username, password) = (%s, %s)", username, password),
		)
	}

	usernameExists, err := dbwrappers.CheckIfUserExists(username)
	if err == nil && usernameExists {
		return "", httperrors.NewHttpError(
			errors.New("username already exists"),
			http.StatusConflict,
			"username already exists",
			currentFilePath+"Signup:usernameExists",
		)
	} else if err != nil {
		return "", httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			currentFilePath+"Signup:usernameExists",
		)
	}

	salt, err := passwords.Salt(username)
	if err != nil {
		return "", httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"could not process credentials",
			currentFilePath+"Signup:password.Salt() with username = (%s)",
		)
	}

	hash := passwords.Hash(password, salt)
	if err = dbwrappers.InsertUser(username, hash, salt); err != nil {
		return "", httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong inserting the new credentials",
			currentFilePath+fmt.Sprintf("Signup:insertUser() with (username, hash, salt) = (%s, %s, %s)", username, hash, salt),
		)
	}

	sessionKey, err := dbwrappers.CreateSessionFor(username)
	if err != nil {
		return "", httperrors.NewHttpError(
			err, //there is no reason to return http error from the createSessionFor func. That layer is the deepest layer in there that starts of the error propagation to the top.
			http.StatusInternalServerError,
			"something went wrong creating the session",
			currentFilePath+fmt.Sprintf("Signup:createSessionFor() with username = %s", username),
		)
	}

	return sessionKey, nil
}
