package handlers

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
)

func SetS3Client(client *s3.Client) {
	s3Client = client
}
