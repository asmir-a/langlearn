package s3wordgame

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/asmir-a/langlearn/backend/aws/s3langlearn"
	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

/*
response needed:
struct response {
	correctWord string
	incorrectWords []string
	correctWordImageUrls []string (to aws s3 bucket)
}

flow:
1) user gets correct word from random
2) we check if there are images for the word in the bucket
	*) we just send the listObjects request to word bucket (this might be somehow avoided i believe, but this is just an optimization)
3) if there are, we return those image urls
4) if not, we send request to google search api
5) we download the images from url to memory and upload them to aws s3 bucket (might fail, right?)
	*) we need to choose image names
		e.g. image-from-web-1.jpg image-from-dalle-1.jpg
6) we send listObjects request
	*) might return less than 4 image urls because download and upload process might fail
		if that happens, we might just retry up to some number of times
7) we get the object names
8) we form the urls and put them to response struct
*/

var S3WordgameBucketName = os.Getenv("S3_WORDGAME_BUCKET_NAME")
var S3WordgameBucketUrlStart = os.Getenv("S3_WORDGAME_BUCKET_URL_START")

func constructS3FileUrl(bucketKey string) string {
	return S3WordgameBucketUrlStart + bucketKey
}

func GetS3FileUrls(word string) ([]string, *httperrors.HttpError) {
	result, err := s3langlearn.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(S3WordgameBucketName),
		Prefix: aws.String(word + "/"),
	})
	if err != nil {
		return nil, httperrors.NewHttp500Error(err)
	}

	fileUrls := []string{}
	for _, content := range result.Contents {
		fileUrls = append(fileUrls, constructS3FileUrl(*content.Key))
	}

	return fileUrls, nil
}

func constructS3WebFileUrlForUpload(word string, taskIndex int, fileRawUrl string) string {
	fileUrl, err := url.Parse(fileRawUrl)
	if err != nil {
		httperrors.Fatal(err)
	}
	extension := path.Ext(fileUrl.Path)

	return word + "/" + "image-from-web-" + strconv.Itoa(taskIndex) + extension
}

func uploadFileToS3(imageKey string, file io.Reader) {
	if _, err := s3langlearn.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(S3WordgameBucketName),
		Key:    aws.String(imageKey),
		Body:   file,
	}); err != nil {
		httperrors.Fatal(err)
	}
}

func uploadFileFromWebToS3(word string, fileUrl string, fileIndexInBucket int) *httperrors.HttpError {
	fileResponse, err := http.Get(fileUrl)
	if err != nil {
		return httperrors.NewHttp500Error(err)
	}
	defer fileResponse.Body.Close()
	fileBytes, err := io.ReadAll(fileResponse.Body)
	if err != nil {
		return httperrors.NewHttp500Error(err)
	}
	fileReader := bytes.NewReader(fileBytes)
	uploadFileToS3(
		constructS3WebFileUrlForUpload(word, fileIndexInBucket, fileUrl),
		fileReader,
	)
	return nil
}

func UploadFilesFromWebToS3(dirName string, fileUrls []string) *httperrors.HttpError {
	numberOfWebImages := len(fileUrls)
	var wg sync.WaitGroup
	var httpErr *httperrors.HttpError
	fmt.Println("starting the goroutines to fetch the images")
	for taskIndex := 1; taskIndex <= numberOfWebImages; taskIndex++ {
		wg.Add(1)
		taskIndex := taskIndex
		go func() {
			defer wg.Done()
			httpErr = uploadFileFromWebToS3(dirName, fileUrls[taskIndex-1], taskIndex)
			fmt.Println(taskIndex, "th routine is done")
		}()
	}
	wg.Wait()
	fmt.Println("all goroutines are done")
	if httpErr != nil {
		return httperrors.WrapError(httpErr)
	}
	return nil
}
