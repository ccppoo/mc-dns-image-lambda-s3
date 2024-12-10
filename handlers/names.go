package handlers

import (
	"fmt"
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
