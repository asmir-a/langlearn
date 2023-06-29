package dbwrappers

import (
	"context"

	"github.com/asmir-a/langlearn/backend/db"
	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/jackc/pgx/v5"
)

func InsertUser(username string, passwordHash string, passwordSalt string) *httperrors.HttpError {
	query := `
		INSERT INTO users (username, password_hash, password_salt)
		VALUES ($1, $2, $3)
	`
	if _, err := db.Conn.Exec(
		context.Background(),
		query,
		username,
		passwordHash,
		passwordSalt,
	); err != nil {
		return httperrors.NewHttp500Error(err)
	}
	return nil
}

func CheckIfUserExists(username string) (bool, *httperrors.HttpError) {
	query := `
		SELECT username
		FROM users
		WHERE username = $1
	`
	var usernameDB string
	err := db.Conn.QueryRow(
		context.Background(),
		query,
		username,
	).Scan(&usernameDB)
	if err == nil {
		return true, nil
	} else if err == pgx.ErrNoRows {
		return false, nil
	} else {
		return false, httperrors.NewHttp500Error(err)
	}
}

func GetUserPasswordHash(username string) (string, *httperrors.HttpError) {
	//assumes that the user with username exists in the database
	query := `
		SELECT password_hash
		FROM users
		WHERE username = $1
	`
	var passwordHashDB string
	if err := db.Conn.QueryRow(
		context.Background(),
		query,
		username,
	).Scan(&passwordHashDB); err != nil {
		return "", httperrors.NewHttp500Error(err)
	}
	return passwordHashDB, nil
}

func GetUserPasswordSalt(username string) (string, *httperrors.HttpError) {
	//assumes that the user with username exists
	query := `
		SELECT password_salt
		FROM users
		WHERE username = $1
	`
	var passwordSaltDB string
	if err := db.Conn.QueryRow(
		context.Background(),
		query,
		username,
	).Scan(&passwordSaltDB); err != nil {
		return "", httperrors.NewHttp500Error(err)
	}
	return passwordSaltDB, nil
}
