package dbwrappers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/asmir-a/langlearn/backend/db"
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
	if _, err := db.Conn.Exec(
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
	_, err := db.Conn.Exec(context.Background(), query, username)
	if err == nil {
		return true, nil
	} else if err == pgx.ErrNoRows {
		return false, nil
	} else {
		return false, httperrors.NewHttp500Error(err)
	}
}

func GetSessionFor(username string) (string, *httperrors.HttpError) {
	query := `
		SELECT session_key
		FROM sessions
		WHERE username = $1
	`
	var sessionKey string
	if err := db.Conn.QueryRow(
		context.Background(),
		query,
		username,
	).Scan(&sessionKey); err == pgx.ErrNoRows { //if session does not exist, it returns pgx.NoRowErr, but it is okay cause the absence of session still means that the session is invalid
		return "", httperrors.NewHttpError(
			err,
			http.StatusUnauthorized,
			"please login first",
		)
	} else if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}
	return sessionKey, nil
}

func DeleteSession(sessionKey string) *httperrors.HttpError {
	query := `
		DELETE FROM sessions
		WHERE session_key = $1
	`
	if _, err := db.Conn.Exec(
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
	if _, err := db.Conn.Exec(
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
	err := db.Conn.QueryRow(
		context.Background(),
		query,
		session_key,
	).Scan(&sessionKeyDb)

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
	if err := db.Conn.QueryRow(
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

func GetUserWith(sessionKey string) (User, *httperrors.HttpError) {
	query := `
		SELECT username
		FROM sessions
		WHERE session_key = $1
	`

	var username string
	if err := db.Conn.QueryRow(
		context.Background(),
		query,
		sessionKey,
	).Scan(&username); err != nil {
		log.Println(err)
		return User{}, httperrors.NewHttp500Error(err)
	}
	//the session also should be checked

	return User{username}, nil //for now, it is enough to use just the username, but in the future, the whole user struct might be more useful
}

//todo: this function has a bug in it. it is triggered when the screen was open for some time and session prolly expired. a bit hard to reproduce.
//scenario is understood now: when the user logs in using another device the session that is stored on the other device gets replaced by a new one. so, when the request for the user is sent with the old cookie, it is not present in the database. this function or the function using this one, should return unauthorized.

//it would be nice if it was possible to wrap error into a new http error
//getuserwith returns an http500 error. however, this is not really useful
//for the client side and also we can return a more meaningful error in this
//case
