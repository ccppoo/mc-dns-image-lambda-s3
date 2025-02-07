package validator

import (
	"image"
)

func isValidateSize(img *image.Image, width int, height int) bool {
	g := (*img).Bounds()

	img_height := g.Dy()
	img_width := g.Dx()
	return img_height == height && img_width == width
}
