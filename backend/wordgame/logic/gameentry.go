package logic

import (
	"encoding/json"

	"github.com/asmir-a/langlearn/backend/httperrors"
	dallewordgame "github.com/asmir-a/langlearn/backend/wordgame/dalle"
	"github.com/asmir-a/langlearn/backend/wordgame/dbwrappers"
	"github.com/asmir-a/langlearn/backend/wordgame/imagesearch"
	s3wordgame "github.com/asmir-a/langlearn/backend/wordgame/s3"
)

type WordGameEntry struct {
	CorrectWord          string   `json:"correctWord"`
	CorrectWordImageUrls []string `json:"correctWordImageUrls"`
	IncorrectWords       []string `json:"incorrectWords"`
}

func getImageUrls(word string, query string) ([]string, *httperrors.HttpError) {
	s3ImageUrls, httpErr := s3wordgame.GetS3FileUrls(word)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}
	if len(s3ImageUrls) > 0 {
		return s3ImageUrls, nil
	}
	webImageUrls, httpErr := imagesearch.FetchImageUrlsFor(query)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}

	dalleImageUrl, httpErr := dallewordgame.GetImageUrlFrom(query)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}

	allImageUrls := append(webImageUrls, dalleImageUrl)
	if httpErr := s3wordgame.UploadFilesFromWebToS3(word, allImageUrls); httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}
	s3ImageUrls, httpErr = s3wordgame.GetS3FileUrls(word)
	if httpErr != nil {
		return nil, httperrors.WrapError(httpErr)
	}
	return s3ImageUrls, nil
}

func calculateMaxIndexForUser(maxIndex int, wordsLearnedCount int) int {
	numberOfLevels := 20
	wordsPerLevel := maxIndex / numberOfLevels

	currentLevel := wordsLearnedCount / wordsPerLevel
	nextLevel := currentLevel + 1

	currentMaxIndexForUser := nextLevel * wordsPerLevel

	if currentMaxIndexForUser < maxIndex {
		return currentMaxIndexForUser
	} else {
		return maxIndex
	}
}

func getGameEntry(username string) (WordGameEntry, *httperrors.HttpError) {
	learnedWordsCount, httpErr := dbwrappers.GetWordsLearnedCount(username)
	if httpErr != nil {
		return WordGameEntry{}, httperrors.WrapError(httpErr)
	}

	maxIndex, httpErr := dbwrappers.GetMaxIndex()

	maxIndexForUser := calculateMaxIndexForUser(maxIndex, learnedWordsCount)
	if httpErr != nil {
		return WordGameEntry{}, httperrors.WrapError(httpErr)
	}

	words, httpErr := dbwrappers.GetRandomKoreanWords(maxIndexForUser)
	if httpErr != nil {
		return WordGameEntry{}, httperrors.WrapError(httpErr)
	}

	correctWordIndex := 0

	correctWord := words[correctWordIndex].Word
	correctWordDef := words[correctWordIndex].Defs[0]

	correctWordImageUrls, httpErr := getImageUrls(correctWord, correctWordDef)
	if httpErr != nil {
		return WordGameEntry{}, httperrors.WrapError(httpErr)
	}

	incorrectWordsWithDefs := words[1:]
	incorrectWords := dbwrappers.ExtractWordsFromWordsWithDefs(incorrectWordsWithDefs) //todo: this does not belong to dbwrappers; also the type prolly should be shared among these packages

	return WordGameEntry{
		CorrectWord:          correctWord,
		CorrectWordImageUrls: correctWordImageUrls,
		IncorrectWords:       incorrectWords,
	}, nil
}

func GetGameEntry(username string) ([]byte, *httperrors.HttpError) {
	gameEntry, httpErr := getGameEntry(username)
	if httpErr != nil {
		return nil, httpErr
	}

	gameEntryBytes, err := json.Marshal(gameEntry)
	if err != nil {
		newHttpErr := httperrors.NewHttp500Error(err)
		return nil, newHttpErr
	}

	return gameEntryBytes, nil
}

/*
/wordgame/users/username/submit
*/
