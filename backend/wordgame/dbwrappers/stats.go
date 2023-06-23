package dbwrappers

import (
	"context"

	"github.com/asmir-a/langlearn/backend/dbconnholder"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

func GetWordsLearnedCount(username string) (int, *httperrors.HttpError) {
	query := `
		SELECT COUNT(*)
		FROM knows
		WHERE username = $1 AND current_count > 3
	`
	var count int
	if err := dbconnholder.Conn.QueryRow(context.Background(), query, username).Scan(&count); err != nil {
		return 0, httperrors.NewHttp500Error(err)
	}
	return count, nil
}

func GetWordsLearningCount(username string) (int, *httperrors.HttpError) {
	query := `
		SELECT COUNT(*)
		FROM knows
		WHERE username = $1
		AND current_count <= 3
	`
	var count int
	if err := dbconnholder.Conn.QueryRow(context.Background(), query, username).Scan(&count); err != nil {
		return 0, httperrors.NewHttp500Error(err)
	}
	return count, nil
}

func GetLearnedWords(username string) {
}

func GetLearningWords(username string) {
}
