package s3langlearn

import (
	"context"

	"github.com/asmir-a/langlearn/backend/httperrors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client

func init() {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		httperrors.Fatal(err)
	}
	S3Client = s3.NewFromConfig(config)
}
