package dbwrappers

import (
	"context"
	"math/rand"
	"time"

	"github.com/asmir-a/langlearn/backend/db"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

const numberOfRandomWords = 4

type WordWithDefs struct {
	Word string
	Defs []string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ExtractWordsFromWordsWithDefs(wordsWithDefs []WordWithDefs) []string {
	words := []string{}
	for _, wordWithDefs := range wordsWithDefs {
		words = append(words, wordWithDefs.Word)
	}
	return words
}

func GetMaxIndex() (int, *httperrors.HttpError) {
	query := `
		SELECT max(index) 
		FROM korean_words
	`
	row := db.Conn.QueryRow(context.Background(), query)
	var maxIndexDb int
	if err := row.Scan(&maxIndexDb); err != nil {
		return 0, httperrors.NewHttp500Error(err)
	}
	return maxIndexDb, nil
}

func GetRandomKoreanWords(maxIndex int) ([]WordWithDefs, *httperrors.HttpError) {
	query := `
		SELECT word, defs
		FROM korean_words
		WHERE index = ANY ($1)
	` //freq info is also prolly needed

	randomIndicesSet := make(map[int]bool)
	for len(randomIndicesSet) != numberOfRandomWords {
		randomIndex := rand.Intn(maxIndex)
		randomIndicesSet[randomIndex] = true
	}

	randomIndices := make([]int, 0)
	for key := range randomIndicesSet {
		randomIndices = append(randomIndices, key)
	}

	wordRows, err := db.Conn.Query(context.Background(), query, randomIndices)
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
