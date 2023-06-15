package dbwrappers

import (
	"context"
	"errors"

	"github.com/asmir-a/langlearn/backend/dbconnholder"
	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/jackc/pgx/v5"
)

//todo: add a constaint for the counter so that it could not be 0; just for safety

type Knows struct {
	Username     string
	Word         string
	CurrentCount int
}

func doesRowExists(username string, word string) (bool, *httperrors.HttpError) {
	query := `
		SELECT username
		FROM knows
		WHERE username = $1 AND word = $2
	`
	var usernameDb string
	err := dbconnholder.Conn.QueryRow(context.Background(), query, username, word).Scan(&usernameDb)
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, httperrors.NewHttp500Error(err)
	}
}

func getCurrentCount(username string, word string) (int, *httperrors.HttpError) {
	query := `
		SELECT current_count
		FROM knows
		WHERE username = $1 AND word = $2
	`
	var currentCount int
	err := dbconnholder.Conn.QueryRow(context.Background(), query, username, word).Scan(&currentCount)
	if err == pgx.ErrNoRows {
		return 0, httperrors.NewHttp500Error(errors.New("expected that the record would exist"))
	} else if err == nil {
		return currentCount, nil
	} else {
		return 0, httperrors.NewHttp500Error(err)
	}
}

func incrementCurrentCount(username string, word string) *httperrors.HttpError {
	query := `
		UPDATE knows
		SET current_count = current_count + 1
		WHERE username = $1 AND word = $2
	`
	if _, err := dbconnholder.Conn.Exec(context.Background(), query, username, word); err != nil {
		return httperrors.NewHttp500Error(err)
	}
	return nil
}

func decrementCurrentCount(username string, word string) *httperrors.HttpError {
	query := `
		UPDATE knows
		SET current_count = current_count - 1
		WHERE username = $1 AND word = $2
	`
	if _, err := dbconnholder.Conn.Exec(context.Background(), query, username, word); err != nil {
		return httperrors.NewHttp500Error(err)
	}
	return nil
}

func createNewRowInKnows(username string, word string) *httperrors.HttpError {
	query := `
		INSERT INTO knows(username, word, word_count) 
		VALUES ($1, $2, $3)
	`
	if _, err := dbconnholder.Conn.Exec(context.Background(), query, username, word, 1); err != nil {
		return httperrors.NewHttp500Error(err)
	}
	return nil
}

func deleteRowInKnows(username string, word string) *httperrors.HttpError {
	query := `
		DELETE FROM knows
		WHERE username = $1 AND word = $2
	`
	if _, err := dbconnholder.Conn.Exec(context.Background(), query, username, word); err != nil {
		return httperrors.NewHttp500Error(err)
	}
	return nil
}
