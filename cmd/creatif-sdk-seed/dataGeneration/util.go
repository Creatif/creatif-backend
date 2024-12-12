package dataGeneration

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	randv2 "math/rand/v2"
	"os"
	"time"
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
		idx := randv2.IntN(len(images))
		pickedImages[i] = images[idx]
	}

	return pickedImages, nil
}

func pickRandomUniqueGroups(groupIds []string, numOfGroups int) []string {
	picked := 0
	groups := make([]string, 3)
	for {
		if picked == numOfGroups {
			return groups
		}

		g := groupIds[randv2.IntN(100)]

		isDuplicate := false
		for _, pickedGroup := range groups {
			if pickedGroup == g {
				isDuplicate = true
			}
		}

		if !isDuplicate {
			groups[picked] = g
			picked++
		}
	}
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

func randomBetween(min, max int) int {
	// Seed the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generate a random number between min and max
	return rand.Intn(max-min+1) + min
}
