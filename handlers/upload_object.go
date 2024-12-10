package handlers

import (
	"log"
	"net/http"

	"strings"

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

	newFileName, fileNameErr := genFileName(fileHeader.Filename)
	contentTypeSplited := strings.Split(fileHeader.Header.Get("Content-Type"), "/")
	if len(contentTypeSplited) < 2 {
		log.Fatalf("invalid file name: missing extension")
	}
	mimeMediaType := contentTypeSplited[0]

	tempFileName := "temp/" + mimeMediaType + "/" + newFileName
	if fileNameErr != nil {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"msg": "malformed file",
		// })
		c.JSON(http.StatusOK, gin.H{
			"success": 0,
			"file": gin.H{
				"url": "",
			},
		})
		return

	}

	bucketName := os.Getenv("BUCKET_PUBLIC")
	CDN_HOST := os.Getenv("CDN_HOST")

	_, err = s3Client.PutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(tempFileName),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		log.Printf("Failed to upload file to S3: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": 0,
			"file": gin.H{
				"url": "",
			},
		})
		return
	}

	imageFileURL := CDN_HOST + "/" + tempFileName
	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"file": gin.H{
			"url": imageFileURL,
		},
	})
}
