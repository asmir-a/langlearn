package imagesearch

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/asmir-a/langlearn/backend/httperrors"
)

var googleImageSearchUrl string = os.Getenv("GOOGLE_IMAGE_SEARCH_URL")
var googleImageSearchEngineId string = os.Getenv("GOOGLE_IMAGE_SEARCH_ENGINE_ID")
var googleImageSearchToken string = os.Getenv("GOOGLE_IMAGE_SEARCH_TOKEN")
var googleImageSearchType string = "image"
var googleImageSearchResultsLimit string = "1"
var googleImageSearchRights string = os.Getenv("GOOGLE_IMAGE_SEARCH_LICENSE")

var googleImageSearchParams map[string]string = map[string]string{
	"cx":         googleImageSearchEngineId,
	"key":        googleImageSearchToken,
	"searchType": googleImageSearchType,
	"rights":     googleImageSearchRights,
	"num":        googleImageSearchResultsLimit,
}

type googleImageSearchResponseItem struct {
	Link string `json:"link"`
}

type googleImageSearchResponseItems struct {
	Items []googleImageSearchResponseItem `json:"items"`
}

func filterAndConvertFromJsonToData(responseBody string) (googleImageSearchResponseItems, *httperrors.HttpError) {
	responseBytes := []byte(responseBody)
	responseData := googleImageSearchResponseItems{}
	if err := json.Unmarshal(responseBytes, &responseData); err != nil {
		return googleImageSearchResponseItems{}, httperrors.NewHttp500Error(err)
	}
	return responseData, nil
}

func fetchResponseItemsFor(query string) (googleImageSearchResponseItems, *httperrors.HttpError) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", googleImageSearchUrl, nil)
	if err != nil {
		return googleImageSearchResponseItems{}, httperrors.NewHttp500Error(err)
	}

	requestUrlParams := url.Values{}
	for k, v := range googleImageSearchParams {
		requestUrlParams.Add(k, v)
	}
	requestUrlParams.Add("q", query)
	req.URL.RawQuery = requestUrlParams.Encode()

	response, err := client.Do(req)
	if err != nil {
		return googleImageSearchResponseItems{}, httperrors.NewHttp500Error(err)
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return googleImageSearchResponseItems{}, httperrors.NewHttp500Error(err)
	}
	bodyString := string(bodyBytes)

	responseData, httpErr := filterAndConvertFromJsonToData(bodyString)
	if httpErr != nil {
		return googleImageSearchResponseItems{}, httperrors.WrapError(httpErr)
	}

	return responseData, nil
}

func FetchImageUrlFor(query string) (string, *httperrors.HttpError) {
	responseItems, httpErr := fetchResponseItemsFor(query)
	if httpErr != nil {
		return "", httperrors.WrapError(httpErr)
	}
	if len(responseItems.Items) == 0 {
		httperrors.Fatal(errors.New("[potential issue: engine quota exceeded]"))
	}
	return responseItems.Items[0].Link, nil
}
