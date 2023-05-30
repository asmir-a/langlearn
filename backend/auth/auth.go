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
	if err = createSession(username); err != nil {
		return err
	}

	return
}

func login(username string, password string) (err error) {
	userExists, err := checkIfUserExists(username)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user does not exist") //maybe need to create a custom error type
	}

	salt := getUserPasswordSalt(username)
	potentialPasswordHash, err := passwords.Hash(password, salt)
	if err != nil {
		return err
	}

	validPasswordHash := getUserPasswordHash(username)
	if potentialPasswordHash != validPasswordHash {
		//todo: may be need to delete the session or not
		return errors.New("wrong credentials")
	}

	sessionExists, err := checkIfSessionExistsFor(username)
	if err != nil {
		return err
	}

	if !sessionExists {
		err = createSessionFor(username)
		return err
	}

	sessionIsValid, err := checkIfSessionIsValidFor(username)
	if !sessionIsValid {
		err = replaceSessionFor(username)
		return err
	}

	err = renewSessionFor(username)
	return err
	//todo: the session should be checked and deleted in a single database transaction
}

func logout(currentSession string) {
	//check if session is correct
	//destroy session
}

func checkIfUserExists(username string) bool {
	_ := getUsername(username) //need to parse the result; maybe, the number of rows or something like that
	return true
}
