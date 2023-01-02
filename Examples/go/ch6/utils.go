package ch6

import (
	"crypto/sha1" //nolint: gosec
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/corona10/goimghdr"
	"gocv.io/x/gocv"
)

const (
	imgPath = "Examples/go/ch6/images"
)

var (
	errUnsupportedMimeType = errors.New("unsupported mim-type")
)

func GetImages() ([]string, error) {
	files, err := os.ReadDir(imgPath)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	var filenames []string

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "jpg") {
			filename := fmt.Sprintf("%s/%s", imgPath, file.Name())

			filenames = append(filenames, filename)
		}
	}

	return filenames, nil
}

func MimeType(filename string) (string, error) {
	mimeType, err := goimghdr.What(filename)
	if err != nil {
		return "", fmt.Errorf("determine image type: %w", err)
	}

	return fmt.Sprintf("image/%s", mimeType), nil
}

func WriteTempFile(obd []byte, mimeType string) (string, error) {
	var filename string

	checksum := sha1.Sum(obd) //nolint: gosec

	switch mimeType {
	case "image/jpg", "image/jpeg":
		filename = fmt.Sprintf("%x.jpg", checksum)
	case "image/png":
		filename = fmt.Sprintf("%x.png", checksum)
	default:
		return "", fmt.Errorf("%w: %s", errUnsupportedMimeType, mimeType)
	}

	filename = fmt.Sprintf("%s/%s", os.TempDir(), filename)
	if err := os.WriteFile(filename, obd, 0600); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return filename, nil
}

func DisplayImage(obd []byte, mimeType string) error {
	filename, err := WriteTempFile(obd, mimeType)
	if err != nil {
		return err
	}

	window := gocv.NewWindow("Image")
	defer window.Close()

	img := gocv.IMRead(filename, gocv.IMReadColor)
	window.IMShow(img)
	window.WaitKey(5000)

	return nil
}
