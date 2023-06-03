package auth

import (
	"errors"

	"github.com/asmir-a/langlearn/backend/auth/passwords"
)

func validateUsername(username string) bool {
	if username == "" {
		return false
	}
	//todo: other checks
	return true
}

func validatePassword(password string) bool {
	if password == "" {
		return false
	}

	return true
}

// assumptions
// there is no session before this begins
func Signup(username string, password string) (string, error) {
	//todo: check if username already exists
	if !validateUsername(username) || !validatePassword(password) {
		return "", errors.New("wrong credentials format")
	}

	salt, err := passwords.Salt(username)
	if err != nil {
		return "", err
	}

	hash, err := passwords.Hash(password, salt)
	if err = insertUser(username, hash, salt); err != nil {
		return "", err
	}

	//todo: need to create a new session for the user
	sessionKey, err := createSessionFor(username)
	if err != nil {
		return "", err
	}

	return sessionKey, nil
}

// todo: need to refactor this function; it is too long
func Login(username string, password string) (session_key string, err error) {
	if !validateUsername(username) || !validatePassword(password) {
		return "", errors.New("wrong credentials")
	}

	userExists, err := checkIfUserExists(username)
	if err != nil {
		return "", err
	}
	if !userExists {
		return "", errors.New("user does not exist") //maybe need to create a custom error type
	}

	salt, err := getUserPasswordSalt(username)
	if err != nil {
		return "", err
	}

	potentialPasswordHash, err := passwords.Hash(password, salt)
	if err != nil {
		return "", err
	}

	validPasswordHash, err := getUserPasswordHash(username)
	if err != nil {
		return "", err
	}

	if potentialPasswordHash != validPasswordHash {
		//todo: may be need to delete the session or not
		return "", errors.New("wrong credentials")
	}

	sessionExists, err := checkIfSessionExistsFor(username) //might be unncessary; we can set up the constraint in the database
	if err != nil {
		return "", err
	}

	if !sessionExists {
		sessionKey, err := createSessionFor(username)
		return sessionKey, err
	}

	sessionKey, err := replaceSessionFor(username) //we do not care if the old session is valid or not; since we got the right login and password, we need to create a new valid session
	return sessionKey, err
	//todo: the session should be checked and deleted in a single database transaction
}

func Logout(currentSessionKey string) error {
	return deleteSession(currentSessionKey)
}
