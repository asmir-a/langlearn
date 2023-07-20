package logic

import (
	"bytes"
	"encoding/json"

	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/dbwrappers"
)

type wordCounts struct {
	Learned  int `json:"learned"`
	Learning int `json:"learning"`
}

func GetWordCountsFor(username string) ([]byte, *httperrors.HttpError) {
	learnedCount, httpErr := dbwrappers.GetWordsLearnedCount(username)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}
	learningCount, httpErr := dbwrappers.GetWordsLearningCount(username)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}
	wordCounts := wordCounts{
		Learned:  learnedCount,
		Learning: learningCount,
	}
	byteBuf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(byteBuf)
	encoder.Encode(wordCounts)
	return byteBuf.Bytes(), nil
}
