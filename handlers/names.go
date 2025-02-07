package handlers

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
)

func genFileName(originalFileName string) (string, error) {

	parts := strings.Split(originalFileName, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid file name: missing extension")
	}
	extension := parts[len(parts)-1] // Get the last part as the extension

	newUUID := strings.ReplaceAll(uuid.New().String(), "-", "")

	// Construct the new file name
	newFileName := fmt.Sprintf("%s.%s", newUUID, extension)
	return newFileName, nil
}

func genTempFileName(fileHeader *multipart.FileHeader) (string, error) {

	newFileName, err := genFileName(fileHeader.Filename)

	if err != nil {
		return "", err
	}

	contentTypeSplited := strings.Split(fileHeader.Header.Get("Content-Type"), "/")

	if len(contentTypeSplited) < 2 {
		return "", fmt.Errorf("invalid file name: missing extension")
	}

	mimeMediaType := contentTypeSplited[0]

	tempFileName := "temp/" + mimeMediaType + "/" + newFileName

	return tempFileName, nil
}
