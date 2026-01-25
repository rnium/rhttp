package images

import (
	"os"
	"path/filepath"
)

const (
	imgNameJpeg = "sample.jpg"
	imgNamePng  = "sample.png"
	imgNameWebp = "sample.webp"
	imgNameSvg  = "sample.svg"
	imgNameGif  = "sample.gif"
)

const (
	acceptPng  = "image/png"
	acceptJpeg = "image/jpeg"
	acceptWebp = "image/webp"
	acceptSvg  = "image/svg+xml"
	acceptGif  = "image/gif"
)

func getImagePath(accept string) string {
	basePath, _ := os.Getwd()
	var imgName string
	switch accept {
	case acceptPng:
		imgName = imgNamePng
	case acceptJpeg:
		imgName = imgNameJpeg
	case acceptWebp:
		imgName = imgNameWebp
	case acceptSvg:
		imgName = imgNameSvg
	case acceptGif:
		imgName = imgNameGif
	default:
		imgName = imgNamePng
	}
	return filepath.Join(basePath, "web/static/images/formats", imgName)
}
