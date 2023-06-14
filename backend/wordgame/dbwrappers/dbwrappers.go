package dbwrappers

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/asmir-a/langlearn/backend/dbconnholder"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

const fileLevelDebugInfo = "wordgame dbwrappers "

const numberOfRandomWords = 4

type WordWithDefs struct {
	Word string
	Defs []string
}

func ExtractWordsFromWordsWithDefs(wordsWithDefs []WordWithDefs) []string {
	words := []string{}
	for _, wordWithDefs := range wordsWithDefs {
		words = append(words, wordWithDefs.Word)
	}
	return words
}

func getMaxIndex() (int, *httperrors.HttpError) {
	const funcLevelDebugInfo = "getMaxIndex "
	query := `
		SELECT max(index) 
		FROM korean_words
	`
	row := dbconnholder.Conn.QueryRow(context.Background(), query)
	var maxIndexDb int
	if err := row.Scan(&maxIndexDb); err != nil {
		httpError := httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			fileLevelDebugInfo+funcLevelDebugInfo+"row.Scan",
		)
		return 0, httpError //todo: we might create an additional layer for errors named DATABASEERROR and pass the error into it. This might be more helpful for debugging. Think about this when you are gonna be refactoring the error structure
	}
	return maxIndexDb, nil
}

func GetRandomKoreanWords() ([]WordWithDefs, *httperrors.HttpError) {
	const funcLevelDebugInfo = "GetRandomKoreanWords "

	query := `
		SELECT word, defs
		FROM korean_words
		WHERE index = ANY ($1)
	` //freq info is also prolly needed

	maxIndex, httpErr := getMaxIndex()
	if httpErr != nil {
		return []WordWithDefs{}, httperrors.WrapError(httpErr, funcLevelDebugInfo+"getMaxIndex()")
	}

	randomIndices := []int{}
	for i := 0; i < numberOfRandomWords; i++ {
		randomIndex := rand.Intn(maxIndex)
		randomIndices = append(randomIndices, randomIndex)
	}

	wordRows, err := dbconnholder.Conn.Query(context.Background(), query, randomIndices)
	if err != nil {
		newHttpErr := httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			fileLevelDebugInfo+funcLevelDebugInfo+"Conn.Query()",
		)
		return []WordWithDefs{}, newHttpErr
	}

	words := []WordWithDefs{}
	for wordRows.Next() {
		word := WordWithDefs{}
		err = wordRows.Scan(&word.Word, &word.Defs)
		if err != nil {
			newHttpErr := httperrors.NewHttpError(
				err,
				http.StatusInternalServerError,
				"something went wrong",
				fileLevelDebugInfo+funcLevelDebugInfo+"wordRows.Scan",
			)
			return []WordWithDefs{}, newHttpErr
		}
		words = append(words, word)
	}

	if wordRows.Err() != nil {
		newHttpErr := httperrors.NewHttpError(
			err,
			http.StatusInternalServerError,
			"something went wrong",
			fileLevelDebugInfo+funcLevelDebugInfo+"wordRows.Err()",
		)
		return []WordWithDefs{}, newHttpErr
	}

	return words, nil
}
