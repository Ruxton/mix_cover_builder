package data

import "image"

type Track struct {
	Artist     string
	Song       string
	Cover      string
	CoverImage image.Image
}
