package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

type Image struct {
	Name      string
	Path      string
	ThumbName string
}

func (i *Image) GenerateThumb(savePath string) {
	file, err := os.Open(i.Path)
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	result := resize.Thumbnail(400, 400, img, resize.NearestNeighbor)

	thumbPath := fmt.Sprintf("%s/%s_thumb.jpg", savePath, i.Name)

	i.AddThumbName(fmt.Sprintf("%s_thumb.jpg", i.Name))

	out, err := os.Create(thumbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	jpeg.Encode(out, result, nil)
}

func (i *Image) AddThumbName(name string) {
	i.ThumbName = name
}
