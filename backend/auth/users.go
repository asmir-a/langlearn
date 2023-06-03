package auth

import (
	"context"

	"github.com/asmir-a/langlearn/backend/database"
	"github.com/jackc/pgx/v5"
)

func insertUser(username string, passwordHash string, passwordSalt string) error {
	query := `
		INSERT INTO users (username, password_hash, password_salt)
		VALUES ($1, $2, $3)
	`
	_, err := database.Conn.Exec(context.Background(), query, username, passwordHash, passwordSalt)
	return err
}

func checkIfUserExists(username string) (bool, error) {
	query := `
		SELECT username
		FROM users
		WHERE username = $1
	`

	var usernameDB string
	err := database.Conn.QueryRow(context.Background(), query, username).Scan(&usernameDB)

	if err == nil {
		return true, nil
	} else if err != nil && err != pgx.ErrNoRows {
		return false, nil
	} else {
		return false, err
	}
}

func getUserPasswordHash(username string) (string, error) {
	//assumes that the user with username exists in the database
	query := `
		SELECT password_hash
		FROM users
		WHERE username = $1
	`
	var passwordHashDB string
	err := database.Conn.QueryRow(context.Background(), query, username).Scan(&passwordHashDB)

	return passwordHashDB, err
}

func getUserPasswordSalt(username string) (string, error) {
	//assumes that the user with username exists
	query := `
		SELECT password_salt
		FROM users
		WHERE username = $1
	`

	var passwordSaltDB string
	err := database.Conn.QueryRow(context.Background(), query, username).Scan(&passwordSaltDB)

	return passwordSaltDB, err
}
