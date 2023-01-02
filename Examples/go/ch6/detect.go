package ch6

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strings"

	"gocv.io/x/gocv"
)

var (
	errLoadCascadeClassifier = errors.New("load cascade classifier from a file")
	errWriteImage            = errors.New("writes a Mat to an image file")
)

func Faces(filename string) (string, error) {
	img := gocv.IMRead(filename, gocv.IMReadGrayScale)

	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	classifierFilename := "Examples/go/ch6/haarcascade_frontalface_alt.xml"
	if !classifier.Load(classifierFilename) {
		return "", fmt.Errorf("%w: %s", errLoadCascadeClassifier, classifierFilename)
	}

	return boxes(
		filename,
		classifier.DetectMultiScaleWithParams(
			img,
			1.2,
			4,
			2,
			image.Point{X: 10, Y: 10},
			image.Point{},
		),
	)
}

func boxes(
	filename string,
	faces []image.Rectangle,
) (string, error) {
	img := gocv.IMRead(filename, gocv.IMReadColor)

	for _, face := range faces {
		gocv.Rectangle(&img, face, color.RGBA{R: 100, G: 100, B: 255}, 2)
	}

	parts := strings.Split(filename, ".")
	newFilename := fmt.Sprintf("%s-detected.%s", parts[0], parts[1])

	if ok := gocv.IMWrite(newFilename, img); !ok {
		return "", fmt.Errorf("%w: %s", errWriteImage, newFilename)
	}

	return newFilename, nil
}
