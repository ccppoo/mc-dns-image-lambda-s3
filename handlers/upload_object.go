package handlers

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	// "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"os"

	"github.com/gin-gonic/gin"
)

func UploadObject(c *gin.Context) {
	var err error

	log.Println("upload obejct")

	// c.Request.Header

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "failed to upload file",
		})
		return
	}

	defer file.Close()

	var output *s3.PutObjectOutput
	bucketName := os.Getenv("BUCKET_STATIC")

	output, err = s3Client.PutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileHeader.Filename),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})

	log.Println(output.ResultMetadata)
	if err != nil {
		log.Printf("Failed to upload receipt to S3: %v", err)

	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Welcome to the players endpoint",
	})
}
