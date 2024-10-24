package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
)

func getListOfImages(dir string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func pickImages(dir string, numOfImages int) ([]os.DirEntry, error) {
	images, err := getListOfImages(dir)
	if err != nil {
		return nil, err
	}

	if len(images) < numOfImages {
		return nil, errors.New("invalid_num_images")
	}

	pickedImages := make([]os.DirEntry, numOfImages)
	for i := 0; i < numOfImages; i++ {
		idx := rand.IntN(len(images))
		pickedImages[i] = images[idx]
	}

	return pickedImages, nil
}

func generateBase64Images(dir string, numOfImages int) ([]string, error) {
	images, err := pickImages(dir, numOfImages)
	if err != nil {
		return nil, err
	}

	b64Images := make([]string, numOfImages)
	for i, img := range images {
		imagePath := fmt.Sprintf("%s/%s", dir, img.Name())
		bytes, err := os.ReadFile(imagePath)
		if err != nil {
			return nil, err
		}

		b64 := fmt.Sprintf("data:image/webp#webp;base64,%s", base64.StdEncoding.EncodeToString(bytes))

		b64Images[i] = b64
	}

	return b64Images, nil
}
