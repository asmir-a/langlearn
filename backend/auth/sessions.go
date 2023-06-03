package auth

import (
	"context"
	"crypto/rand"
	"log"
	"time"

	"github.com/asmir-a/langlearn/backend/database"
	"github.com/jackc/pgx/v5"
)

const SESSION_ID_LENGTH = 128
const SESSION_LOGIN_DURATION = time.Duration(time.Hour * 24 * 7)
const SESSION_USAGE_DURATION = time.Duration(time.Hour * 24)

func generateSessionKey(username string) string {
	randomBytes := make([]byte, SESSION_ID_LENGTH)
	if _, err := rand.Read(randomBytes); err != nil { //todo: think more if this is secure
		log.Fatal("could not generate random bytes: ", err)
	}
	return username + string(randomBytes)
}

func createSessionFor(username string) (string, error) {
	query := `
		INSERT INTO sessions (session_key, username, login_time, last_seen_time)
		VALUES ($1, $2, $3, $4)
	`
	newSessionKey := generateSessionKey(username)
	loginTime := time.Now()
	lastSeenTime := time.Now()

	_, err := database.Conn.Exec(context.Background(), query, newSessionKey, username, loginTime, lastSeenTime)
	return newSessionKey, err
}

func replaceSessionFor(username string) (string, error) {
	err := deleteSessionFor(username)
	if err != nil {
		return "", err
	}

	sessionKey, err := createSessionFor(username)
	return sessionKey, err
}

func checkIfSessionExistsFor(username string) (bool, error) {
	query := `
		SELECT * FROM sessions
		WHERE username = $1
	`
	_, err := database.Conn.Exec(context.Background(), query, username)

	if err == nil {
		return true, nil
	} else if err == pgx.ErrNoRows {
		return false, nil
	} else {
		return false, err
	}
}

// for now, anytime you update the session, it is better to replace it as a whole cause the login might be coming from another device
// with that in mind, each user prolly needs to have exactly one session in the database
// nope, the user can still logout. in that case, the database should delete the session entry
func checkIfSessionIsValidFor(username string) (bool, error) {
	query := `
		SELECT session_key, username, login_time, last_seen_time
		FROM sessions
		WHERE username = $1
	`
	var sessionKeyDB, usernameDB string
	var loginTimeDB, lastSeenTimeDB time.Time
	err := database.Conn.QueryRow(
		context.Background(),
		query,
		username,
	).Scan(
		&sessionKeyDB,
		&usernameDB,
		&loginTimeDB,
		&lastSeenTimeDB,
	)
	if err != nil {
		return false, err
	}

	loggedInTooLongAgo := time.Now().Sub(loginTimeDB) > SESSION_LOGIN_DURATION
	lastSeenTooLongAgo := time.Now().Sub(lastSeenTimeDB) > SESSION_USAGE_DURATION
	if loggedInTooLongAgo || lastSeenTooLongAgo {
		return false, nil
	} else {
		return true, nil
	}
}

func deleteSession(sessionKey string) error {
	query := `
		DELETE FROM sessions
		WHERE session_key = $1
	`
	_, err := database.Conn.Exec(context.Background(), query, sessionKey)
	return err
}

func deleteSessionFor(username string) error {
	query := `
		DELETE FROM sessions
		WHERE username = $1
	`

	_, err := database.Conn.Exec(context.Background(), query, username)
	if err != nil {
		return err
	}

	return nil
}

func checkIfSessionExists(session_key string) (bool, error) {
	query := `
		SELECT *
		FROM sessions
		WHERE session_key = $1
	`
	_, err := database.Conn.Exec(context.Background(), query, session_key)
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func CheckIfSessionIsValid(session_key string) (bool, error) {
	sessionExists, err := checkIfSessionExists(session_key)
	if err != nil {
		return false, err
	}
	if !sessionExists { //todo: this might even be impossible; can just put assert(0) or leave it empty
		return false, nil
	}

	query := `
		SELECT login_time, last_seen_time
		FROM sessions
		WHERE session_key = $1
	`
	var loginTime, lastSeenTime time.Time
	err = database.Conn.QueryRow(context.Background(), query, session_key).Scan(&loginTime, &lastSeenTime)
	if err != nil {
		return false, err
	}
	loggedInTooLongAgo := time.Now().Sub(loginTime) > SESSION_LOGIN_DURATION
	lastSeenTooLongAgo := time.Now().Sub(lastSeenTime) > SESSION_USAGE_DURATION

	if loggedInTooLongAgo || lastSeenTooLongAgo {
		return false, nil
	}

	return true, nil
}

//todo: for the next project, try to use a migration tool like goose or go-migrate.
