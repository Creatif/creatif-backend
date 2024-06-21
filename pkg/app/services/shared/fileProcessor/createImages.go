package fileProcessor

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/google/uuid"
	"os"
	"regexp"
	"strings"
)

type tempImage struct {
	path        string
	base64Image *string
}

type createdFile struct {
	path     string
	filePath string
}

type uploadResult struct {
	createdFile createdFile
	error       error
}

func UploadFiles(projectId string, value []byte, imagePaths []string) ([]byte, error) {
	jsonParsed, err := gabs.ParseJSON(value)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Parsing JSON failed: %s", err.Error()))
	}

	base64Images := make([]tempImage, 0)
	for _, path := range imagePaths {
		base64Image, ok := jsonParsed.Path(path).Data().(string)
		if !ok {
			return nil, errors.New(fmt.Sprintf("Could not find path: %s", path))
		}

		base64Images = append(base64Images, tempImage{
			path:        path,
			base64Image: &base64Image,
		})

		_, err := jsonParsed.Set(nil, path)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not nullify path: %s", path))
		}
	}

	createdFiles, err := uploadFiles(projectId, base64Images)
	fmt.Println(err)
	if err != nil {
		for _, file := range createdFiles {
			os.Remove(file.filePath)
		}

		return nil, err
	}

	fmt.Println(createdFiles)

	for _, file := range createdFiles {
		_, err := jsonParsed.Set(file.filePath, file.path)
		if err != nil {
			for _, file := range createdFiles {
				os.Remove(file.filePath)
			}

			return value, err
		}
	}

	return jsonParsed.Bytes(), nil
}

func uploadFiles(projectId string, base64Images []tempImage) ([]createdFile, error) {
	createdFiles := make([]createdFile, 0)

	for _, image := range base64Images {
		spl := strings.Split(*image.base64Image, "base64,")
		dec, err := base64.StdEncoding.DecodeString(spl[1])
		if err != nil {
			return createdFiles, fmt.Errorf("Could not decode base64 image: %w", err)
		}

		extension, err := extractAndValidateMimeType(image.base64Image)
		if err != nil {
			return createdFiles, err
		}

		fileName := fmt.Sprintf("%s.%s", uuid.NewString(), extension)
		filePath := fmt.Sprintf(
			"%s/%s/%s",
			os.Getenv("PUBLIC_DIRECTORY"),
			projectId,
			fileName,
		)

		f, err := os.Create(filePath)
		if err != nil {
			return createdFiles, err
		}

		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			return createdFiles, err
		}

		if err := f.Sync(); err != nil {
			return createdFiles, err
		}

		createdFiles = append(createdFiles, createdFile{
			path:     image.path,
			filePath: fmt.Sprintf("/api/v1/static/%s/%s", projectId, fileName),
		})
	}

	return createdFiles, nil
}

func extractAndValidateMimeType(image *string) (string, error) {
	re := regexp.MustCompile(`data:(.*);`)
	match := re.FindStringSubmatch(*image)
	if match == nil {
		return "", errors.New("Could not determine mime type")
	}

	if len(match) < 2 {
		return "", errors.New("Could not determine mime type")
	}

	mimeType := match[1]
	sep := strings.Split(mimeType, "#")
	if len(sep) < 2 {
		return "", errors.New("base64 has mime type but could not determine extension")
	}

	return sep[1], nil
}
