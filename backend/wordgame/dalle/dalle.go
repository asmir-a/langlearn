package dallewordgame

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/asmir-a/langlearn/backend/httperrors"
)

var openAiDalleUrl = os.Getenv("OPENAI_DALLE_URL")
var openAiDalleKey = os.Getenv("OPENAI_DALLE_KEY")
var openAiDalleN = 1
var openAiDalleSize = "512x512"

type dalleRequestParams struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type urlEntry struct {
	Url string `json:"url"`
}

type dalleResponse struct {
	Created int        `json:"created"`
	Data    []urlEntry `json:"data"`
}

func init() {
	log.Println("the dalle token is: ", openAiDalleKey)
}

func turnPromptToMoreExpressive(prompt string) string {
	return fmt.Sprintf("an illustrative painting of the concept of %s with high contrasting colors", prompt)
}

func GetImageUrlFrom(prompt string) (string, *httperrors.HttpError) {
	dalleRequestBody := dalleRequestParams{
		Prompt: turnPromptToMoreExpressive(prompt),
		N:      openAiDalleN,
		Size:   openAiDalleSize,
	}
	dalleRequestBodyJson, err := json.Marshal(dalleRequestBody)
	if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}
	dalleRequestBuffer := bytes.NewBuffer(dalleRequestBodyJson)

	clientForDalle := http.Client{}
	requestToDalle, err := http.NewRequest("POST", openAiDalleUrl, dalleRequestBuffer)
	if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}
	requestToDalle.Header.Set("Content-Type", "application/json")
	requestToDalle.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openAiDalleKey))

	responseFromDalle, err := clientForDalle.Do(requestToDalle)
	if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}
	defer responseFromDalle.Body.Close()

	responseBodyJson, err := io.ReadAll(responseFromDalle.Body)
	if err != nil {
		return "", httperrors.NewHttp500Error(err)
	}
	var responseBody dalleResponse
	json.Unmarshal(responseBodyJson, &responseBody)

	return responseBody.Data[0].Url, nil
}
