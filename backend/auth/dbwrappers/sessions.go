package dbwrappers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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
	if _, err := dbconnholder.Conn.Exec(
		context.Background(),
		query,
		newSessionKey,
		username,
		loginTime,
		lastSeenTime,
	); err != nil {
		return "", httperrors.NewHttp500Error(err)
	}
	return newSessionKey, nil
}

func ReplaceSessionFor(username string) (string, *httperrors.HttpError) {
	if httpErr := DeleteSessionFor(username); httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	}
	if sessionKey, httpErr := CreateSessionFor(username); httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	} else {
		return sessionKey, nil
	}
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

func DeleteSession(sessionKey string) *httperrors.HttpError {
	query := `
		DELETE FROM sessions
		WHERE session_key = $1
	`
	if _, err := dbconnholder.Conn.Exec(
		context.Background(),
		query,
		sessionKey,
	); err != nil {
		return httperrors.NewHttp500Error(err)
	}
	return nil
}

func DeleteSessionFor(username string) *httperrors.HttpError {
	query := `
		DELETE FROM sessions
		WHERE username = $1
	`
	if _, err := dbconnholder.Conn.Exec(
		context.Background(),
		query,
		username,
	); err != nil {
		return httperrors.NewHttp500Error(err)
	}
	return nil
}

func checkIfSessionExists(session_key string) (bool, *httperrors.HttpError) {
	query := `
		SELECT session_key
		FROM sessions
		WHERE session_key = $1
	`
	var sessionKeyDb string
	err := dbconnholder.Conn.QueryRow(context.Background(), query, session_key).Scan(&sessionKeyDb)

	if err == nil {
		return true, nil
	} else if err == pgx.ErrNoRows {
		return false, nil
	} else {
		return false, httperrors.NewHttp500Error(err)
	}
}

func CheckIfSessionIsValid(sessionKey string) (bool, *httperrors.HttpError) {
	sessionExists, httpErr := checkIfSessionExists(sessionKey)
	if httpErr != nil {
		return false, httperrors.WrapError(httpErr)
	} else if !sessionExists { //todo: this might even be impossible; can just put assert(0) or leave it empty
		return false, nil
	}

	query := `
		SELECT login_time, last_seen_time
		FROM sessions
		WHERE session_key = $1
	`
	var loginTime, lastSeenTime time.Time
	if err := dbconnholder.Conn.QueryRow(
		context.Background(),
		query,
		sessionKey,
	).Scan(&loginTime, &lastSeenTime); err != nil {
		return false, httperrors.NewHttp500Error(err)
	}

	loggedInTooLongAgo := time.Now().Sub(loginTime) > SESSION_LOGIN_DURATION
	lastSeenTooLongAgo := time.Now().Sub(lastSeenTime) > SESSION_USAGE_DURATION
	if loggedInTooLongAgo || lastSeenTooLongAgo {
		return false, nil
	}
	return true, nil
}

func GetUserWith(session_key string) (*User, *httperrors.HttpError) {
	query := `
		SELECT username
		FROM sessions
		where session_key = $1
	`

	var username string
	if err := dbconnholder.Conn.QueryRow(context.Background(), query, session_key).Scan(&username); err != nil {
		return nil, httperrors.NewHttp500Error(err)
	}

	return &User{username}, nil //for now, it is enough to use just the username, but in the future, the whole user struct might be more useful
}
