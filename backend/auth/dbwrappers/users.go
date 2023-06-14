package dbwrappers

import (
	"context"

	"github.com/asmir-a/langlearn/backend/dbconnholder"
	"github.com/jackc/pgx/v5"
)

func InsertUser(username string, passwordHash string, passwordSalt string) error {
	query := `
		INSERT INTO users (username, password_hash, password_salt)
		VALUES ($1, $2, $3)
	`
	_, err := dbconnholder.Conn.Exec(context.Background(), query, username, passwordHash, passwordSalt)
	return err
}

func CheckIfUserExists(username string) (bool, error) {
	query := `
		SELECT username
		FROM users
		WHERE username = $1
	`

	var usernameDB string
	err := dbconnholder.Conn.QueryRow(context.Background(), query, username).Scan(&usernameDB)

	if err == nil {
		return true, nil
	} else if err != nil && err != pgx.ErrNoRows {
		return false, nil
	} else {
		return false, err
	}
}

func GetUserPasswordHash(username string) (string, error) {
	//assumes that the user with username exists in the database
	query := `
		SELECT password_hash
		FROM users
		WHERE username = $1
	`
	var passwordHashDB string
	err := dbconnholder.Conn.QueryRow(context.Background(), query, username).Scan(&passwordHashDB)

	return passwordHashDB, err
}

func GetUserPasswordSalt(username string) (string, error) {
	//assumes that the user with username exists
	query := `
		SELECT password_salt
		FROM users
		WHERE username = $1
	`

	var passwordSaltDB string
	err := dbconnholder.Conn.QueryRow(context.Background(), query, username).Scan(&passwordSaltDB)

	return passwordSaltDB, err
}
