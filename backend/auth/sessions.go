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

func createOrRenewSession(username string) (err error) {
	query := `
		INSERT INTO sessions (session_key, username, login_time, last_seen_time)
		VALUES ($1, $2, $3, $4)
	`
	database.Conn.Begin(context.Background())
}

func generateSessionKey(username string) string {
	randomBytes := make([]byte, SESSION_ID_LENGTH)
	if _, err := rand.Read(randomBytes); err != nil { //todo: think more if this is secure
		log.Fatal("could not generate random bytes: ", err)
	}
	return username + string(randomBytes)
}

func checkIfSessionExists(username string) (bool, error) {
	query := `
		SELECT * FROM sessions
		WHERE username = $1
	`
	_, err := database.Conn.Exec(context.Background(), query, username)
	if err == nil {
		return true, nil
	}
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return false, err
}

func isSessionValid(username string) (bool, error) {
	query := `
		SELECT * FROM sessions
		WHERE username = $1
	`

	var sessionKeyDb, usernameDb string
	var loginTimeDb, lastSeenTimeDb time.Time

	err := database.Conn.QueryRow(context.Background(), query, username).Scan(&sessionKeyDb, &usernameDb, &loginTimeDb, &lastSeenTimeDb)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Fatal("session must be valid") //in the future transactions will be used to handle this case
		}
		return false, err
	}

	if time.Now().Sub(loginTimeDb) > SESSION_LOGIN_DURATION {
		return false, nil
	}

	if time.Now().Sub(lastSeenTimeDb) > SESSION_USAGE_DURATION {
		return false, nil
	}
}

func checkIfSessionExistsAndValid(username string) (bool, error) {
}

func deleteSession(username string) error {
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

func replaceSession(username string) error {
}

func createSession(username string) (err error) {
	//for now, we are sure that there is no session in the table because this is called by signup
	query := `
		INSERT INTO sessions (session_key, username, login_time, last_seen_time)
		VALUES ($1, $2, $3, $4)
	`

	_, err = database.Conn.Exec(
		context.Background(),
		query,
		generateSessionKey(username),
		username,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		log.Fatal("error generating a session key: ", err)
	}

	return
}

func destroySession(session_key string) {
}

func renewSessionLoginTime() {
}

func renewSessionLastSeenTime() {
}

//todo: for the next project, try to use a migration tool like goose or go-migrate.
