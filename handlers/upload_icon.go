package handlers

import (
	"log"
	"net/http"

	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	// "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"bytes"
	"image"
	_ "image/png"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func UploadIcon(c *gin.Context) {
	var err error

	log.Println("upload icon")

	// file, fileHeader, err := c.Request.FormFile("file")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}

	log.Println("open image file")

	openedFile, err := file.Open()
	if err != nil {
		// log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"file": gin.H{
				"url": "",
			},
		})
		return
	}
	defer openedFile.Close()

	log.Println("tee reader")

	var buf bytes.Buffer
	tee := io.TeeReader(openedFile, &buf)

	m, _, err := image.Decode(tee)
	// png 형식이 아닐 경우 error
	log.Println("tee reader")

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"file": gin.H{
				"url": "",
			},
		})
		return
	}

	log.Println("get image size")

	g := m.Bounds()

	// Get height and width
	height := g.Dy()
	width := g.Dx()
	log.Printf("image width : %d, height : %d", width, height)

	if !(height == 64 && width == 64) {
		// log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"file": gin.H{
				"url": "",
			},
		})
		return
	}

	newFileName, fileNameErr := genFileName(file.Filename)
	contentTypeSplited := strings.Split(file.Header.Get("Content-Type"), "/")
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
		Body:        &buf,
		ContentType: aws.String(file.Header.Get("Content-Type")),
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
