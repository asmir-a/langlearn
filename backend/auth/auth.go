package auth

import (
	"errors"

	"github.com/asmir-a/langlearn/backend/auth/passwords"
)

// assumptions
// there is no session before this begins
func signup(username string, password string) (err error) {
	//todo: check if username already exists
	salt, err := passwords.Salt(username)
	if err != nil {
		return err
	}

	hash, err := passwords.Hash(password, salt)
	if err = insertUser(username, hash, salt); err != nil {
		return err
	}

	//todo: need to create a new session for the user
	if err = createSessionFor(username); err != nil {
		return err
	}

	return
}

// todo: need to refactor this function; it is too long
func login(username string, password string) (err error) {
	userExists, err := checkIfUserExists(username)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user does not exist") //maybe need to create a custom error type
	}

	salt, err := getUserPasswordSalt(username)
	if err != nil {
		return err
	}

	potentialPasswordHash, err := passwords.Hash(password, salt)
	if err != nil {
		return err
	}

	validPasswordHash, err := getUserPasswordHash(username)
	if err != nil {
		return err
	}

	if potentialPasswordHash != validPasswordHash {
		//todo: may be need to delete the session or not
		return errors.New("wrong credentials")
	}

	sessionExists, err := checkIfSessionExistsFor(username) //might be unncessary; we can set up the constraint in the database
	if err != nil {
		return err
	}

	if !sessionExists {
		err = createSessionFor(username)
		return err
	}

	err = replaceSessionFor(username) //we do not care if the old session is valid or not; since we got the right login and password, we need to create a new valid session
	return err
	//todo: the session should be checked and deleted in a single database transaction
}

func logout(currentSessionKey string) error {
	//check if session is correct
	//destroy session
	return deleteSession(currentSessionKey)
}
