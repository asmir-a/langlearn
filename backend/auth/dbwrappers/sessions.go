package dbwrappers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/asmir-a/langlearn/backend/dbconnholder"
	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/jackc/pgx/v5"
)

const SESSION_ID_LENGTH = 128
const SESSION_LOGIN_DURATION = time.Duration(time.Hour * 24 * 7)
const SESSION_USAGE_DURATION = time.Duration(time.Hour * 24)

func generateSessionKey(username string) string {
	randomBytes := make([]byte, SESSION_ID_LENGTH)
	if _, err := rand.Read(randomBytes); err != nil { //todo: think more if this is secure
		log.Fatal("should be impossible")
	}
	randomString := base64.StdEncoding.EncodeToString(randomBytes)
	return username + randomString
}

func CreateSessionFor(username string) (string, *httperrors.HttpError) {
	query := `
		INSERT INTO sessions (session_key, username, login_time, last_seen_time)
		VALUES ($1, $2, $3, $4)
	`
	newSessionKey := generateSessionKey(username)
	loginTime := time.Now()
	lastSeenTime := time.Now()

	_, err := dbconnholder.Conn.Exec(context.Background(), query, newSessionKey, username, loginTime, lastSeenTime)
	return newSessionKey, httperrors.NewHttp500Error(err)
}

func ReplaceSessionFor(username string) (string, *httperrors.HttpError) {
	err := DeleteSessionFor(username)
	if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}

	sessionKey, err := CreateSessionFor(username)
	return sessionKey, httperrors.NewHttp500Error(err)
}

func CheckIfSessionExistsFor(username string) (bool, *httperrors.HttpError) {
	query := `
		SELECT * FROM sessions
		WHERE username = $1
	`
	_, err := dbconnholder.Conn.Exec(context.Background(), query, username)

	if err == nil {
		return true, nil
	} else if err == pgx.ErrNoRows {
		return false, nil
	} else {
		return false, httperrors.NewHttp500Error(err)
	}
}

// for now, anytime you update the session, it is better to replace it as a whole cause the login might be coming from another device
// with that in mind, each user prolly needs to have exactly one session in the database
// nope, the user can still logout. in that case, the database should delete the session entry
func CheckIfSessionIsValidFor(username string) (bool, *httperrors.HttpError) {
	query := `
		SELECT session_key, username, login_time, last_seen_time
		FROM sessions
		WHERE username = $1
	`
	var sessionKeyDB, usernameDB string
	var loginTimeDB, lastSeenTimeDB time.Time
	err := dbconnholder.Conn.QueryRow(
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
		return false, httperrors.NewHttp500Error(err)
	}

	loggedInTooLongAgo := time.Now().Sub(loginTimeDB) > SESSION_LOGIN_DURATION
	lastSeenTooLongAgo := time.Now().Sub(lastSeenTimeDB) > SESSION_USAGE_DURATION
	if loggedInTooLongAgo || lastSeenTooLongAgo {
		return false, nil
	} else {
		return true, nil
	}
}

func DeleteSession(sessionKey string) *httperrors.HttpError {
	query := `
		DELETE FROM sessions
		WHERE session_key = $1
	`
	_, err := dbconnholder.Conn.Exec(context.Background(), query, sessionKey)
	return httperrors.NewHttp500Error(err)
}

func DeleteSessionFor(username string) *httperrors.HttpError {
	query := `
		DELETE FROM sessions
		WHERE username = $1
	`

	_, err := dbconnholder.Conn.Exec(context.Background(), query, username)
	if err != nil {
		return httperrors.NewHttp500Error(err)
	}

	return nil
}

func CheckIfSessionExists(session_key string) (bool, *httperrors.HttpError) {
	query := `
		SELECT session_key
		FROM sessions
		WHERE session_key = $1
	`
	row := dbconnholder.Conn.QueryRow(context.Background(), query, session_key)
	//todo!:fix these bad things
	var session_keyFromDb string
	err := row.Scan(&session_keyFromDb)

	if err == pgx.ErrNoRows {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, httperrors.NewHttp500Error(err)
	}
}

func CheckIfSessionIsValid(session_key string) (bool, *httperrors.HttpError) {
	sessionExists, httpErr := CheckIfSessionExists(session_key)
	if httpErr != nil {
		return false, httpErr
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
	err := dbconnholder.Conn.QueryRow(context.Background(), query, session_key).Scan(&loginTime, &lastSeenTime)
	if err != nil {
		fmt.Println("query row error")
		fmt.Println("error is: ", err)
		return false, httperrors.NewHttp500Error(err)
	}
	loggedInTooLongAgo := time.Now().Sub(loginTime) > SESSION_LOGIN_DURATION
	lastSeenTooLongAgo := time.Now().Sub(lastSeenTime) > SESSION_USAGE_DURATION

	if loggedInTooLongAgo || lastSeenTooLongAgo {
		return false, nil
	}

	return true, nil
}
