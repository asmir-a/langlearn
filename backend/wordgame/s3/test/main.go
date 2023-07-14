package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//database:
/*

wordsToImagesTable
(one-to-many)
[wordId][imageUrl]

(imageUrl: link to s3)

fetch the urls from the s3
if the number of items is less than 4
need to populate the buckets once more

if it is 4, just return the values

we can make the bucket public. or we can use cdn to store the images.(storing the images in cdn would not trigger const increase after a certain point cause there is a fixed number of words in the dictionary)
we can serve the images from the server.
*/

var s3Client *s3.Client
var bucketName = "langlearn-words"

func AddUrlToDb(word string, bucketUrl string, client *s3.Client) error {
	return nil
}

func TryToGetImageUrls(word string, s3Client *s3.Client) {
	result, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(word),
	})
	if err != nil {
		fmt.Println("something went wrong: ", err)
		panic(err)
	}
	fmt.Println("the contents of the bucket are: ", len(result.Contents))
	for _, content := range result.Contents {
		fmt.Println("the key of content is: ", *content.Key)
	}
}

func GetImageViaUrl(imageUrl string) error {
	response, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println("there was an error: ", err)
		return err
	}
	defer response.Body.Close()

	newFile, err := os.Create("image.jpg")
	if err != nil {
		fmt.Println("something went wrong: ", err)
		return err
	}
	io.Copy(newFile, response.Body)
	return nil
}

func DownloadAndUploadImage(word string, imageUrl string, client *s3.Client) error {
	imageResponse, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println("get error: ", err)
		return err
	}
	fmt.Println("response status: ", imageResponse.Status)
	defer imageResponse.Body.Close()

	imageFile, err := io.ReadAll(imageResponse.Body)
	if err != nil {
		fmt.Println("error saving the file into memory: ", err)
		return err
	}
	imageFileBuffer := bytes.NewBuffer(imageFile)

	extension := filepath.Ext(imageUrl)
	s3FullFilePath := word + "/" + "image-from-web" + extension

	output, err := client.PutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String("langlearn-words"),
			Key:    aws.String(s3FullFilePath),
			Body:   imageFileBuffer,
		})

	if err != nil {
		fmt.Println("s3 error: ", err)
		return err
	}

	fmt.Println("s3 put object output: ", output)

	return nil
}

func ListObjects(client *s3.Client) {
	output, err := client.ListObjects(context.TODO(), &s3.ListObjectsInput{
		Bucket: aws.String("asmir-test-bucket"),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("first page results: ")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
}

var testWord string = "가지"
var testImageUrl string = "https://upload.wikimedia.org/wikipedia/commons/thumb/f/fb/Full_Moon_-_Flickr_-_gjdonatiello.jpg/320px-Full_Moon_-_Flickr_-_gjdonatiello.jpg"
var s3ImageUrl string = fmt.Sprintf("http://langlearn-words.s3.ap-northeast-2.amazonaws.com/%s/image-from-web.jpg", testWord)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(cfg)

	startTime := time.Now().UnixMilli()
	DownloadAndUploadImage(testWord, testImageUrl, client)
	endTime := time.Now().UnixMilli()

	TryToGetImageUrls(testWord, client)
	GetImageViaUrl(s3ImageUrl)

	duration := endTime - startTime
	fmt.Printf("it took %d milliseconds to upload and download the file\n", duration)
}
