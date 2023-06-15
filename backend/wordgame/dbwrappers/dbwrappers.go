package dbwrappers

import (
	"context"
	"math/rand"

	"github.com/asmir-a/langlearn/backend/dbconnholder"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

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
	query := `
		SELECT max(index) 
		FROM korean_words
	`
	row := dbconnholder.Conn.QueryRow(context.Background(), query)
	var maxIndexDb int
	if err := row.Scan(&maxIndexDb); err != nil {
		return 0, httperrors.NewHttp500Error(err)
	}
	return maxIndexDb, nil
}

func GetRandomKoreanWords() ([]WordWithDefs, *httperrors.HttpError) {
	query := `
		SELECT word, defs
		FROM korean_words
		WHERE index = ANY ($1)
	` //freq info is also prolly needed

	maxIndex, httpErr := getMaxIndex()
	if httpErr != nil {
		return []WordWithDefs{}, httperrors.WrapError(httpErr)
	}

	randomIndices := []int{}
	for i := 0; i < numberOfRandomWords; i++ {
		randomIndex := rand.Intn(maxIndex)
		randomIndices = append(randomIndices, randomIndex)
	}

	wordRows, err := dbconnholder.Conn.Query(context.Background(), query, randomIndices)
	if err != nil {
		return []WordWithDefs{}, httperrors.NewHttp500Error(err)
	}

	words := []WordWithDefs{}
	for wordRows.Next() {
		word := WordWithDefs{}
		err = wordRows.Scan(&word.Word, &word.Defs)
		if err != nil {
			return []WordWithDefs{}, httperrors.NewHttp500Error(err)
		}
		words = append(words, word)
	}

	if wordRows.Err() != nil {
		return []WordWithDefs{}, httperrors.NewHttp500Error(err)
	}

	return words, nil
}
