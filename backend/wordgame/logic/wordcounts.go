package logic

import (
	"bytes"
	"encoding/json"

	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/asmir-a/langlearn/backend/wordgame/dbwrappers"
)

type WordCounts struct {
	LearnedCount  int `json:"learnedCount"`
	LearningCount int `json:"learningCount"`
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
	wordCounts := WordCounts{
		LearnedCount:  learnedCount,
		LearningCount: learningCount,
	}
	byteBuf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(byteBuf)
	encoder.Encode(wordCounts)
	return byteBuf.Bytes(), nil
}
