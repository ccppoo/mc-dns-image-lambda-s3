package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/gin-gonic/gin"
)

func ListObjects(c *gin.Context) {
	var err error

	var output *s3.ListObjectsV2Output
	var objects []types.Object

	bucketName := os.Getenv("BUCKET_STATIC")
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}
	objectPaginator := s3.NewListObjectsV2Paginator(s3Client, input)

	for objectPaginator.HasMorePages() {
		output, err = objectPaginator.NextPage(c.Request.Context())
		if err != nil {
			var noBucket *types.NoSuchBucket
			if errors.As(err, &noBucket) {
				log.Printf("Bucket %s does not exist.\n", bucketName)
				err = noBucket
			}
			break
		} else {
			objects = append(objects, output.Contents...)
		}
	}
	for i, object := range objects {
		fmt.Printf("%d : %s, size : %d\n", i, *object.Key, *object.Size)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Welcome to the players endpoint",
	})
}
