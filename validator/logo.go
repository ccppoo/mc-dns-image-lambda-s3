package validator

import (
	"log"

	"bytes"
	"fmt"
	"image"
	png "image/png"
	"io"
)

// Logo contraints
// 1. png
// 2. size 64x64
func SanitizeLogo(reader io.Reader) (bytes.Buffer, error) {
	img, format, err := image.Decode(reader)
	var buf bytes.Buffer

	if err != nil {
		log.Fatal(err)
		return buf, fmt.Errorf("failed to decode image")
	}

	if format != "png" {
		return buf, fmt.Errorf("file is not PNG format")
	}

	if !isValidateSize(&img, 64, 64) {
		return buf, fmt.Errorf("not correct size")
	}

	err = reEncodePNG(&buf, &img)

	if err != nil {
		return buf, err
	}

	return buf, nil
}

func reEncodePNG(buffer *bytes.Buffer, img *image.Image) error {
	err := png.Encode(buffer, *img)
	if err != nil {
		return fmt.Errorf("failed to encode png file")
	}
	return nil
}
