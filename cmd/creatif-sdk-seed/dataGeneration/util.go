package dataGeneration

import (
	"errors"
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

func randomBetween(min, max int) int {
	// Seed the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generate a random number between min and max
	return rand.Intn(max-min+1) + min
}
