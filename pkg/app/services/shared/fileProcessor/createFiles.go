package fileProcessor

import (
	"creatif/pkg/app/services/events"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"os"
)

func UploadFiles(projectId string, value []byte, imagePaths []string, callback callbackCreateFn) ([]byte, error) {
	jsonParsed, err := gabs.ParseJSON(value)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Parsing JSON failed: %s", err.Error()))
	}

	uploadedFiles := make([]createdFile, 0)
	var processingError error

	defer func() {
		if processingError == nil {
			return
		}

		for _, file := range uploadedFiles {
			if err := os.Remove(file.FileSystemFilePath); err != nil {
				events.DispatchEvent(events.NewFileNotRemoveEvent(file.FileSystemFilePath, "", projectId))
			}
		}
	}()

	for _, path := range imagePaths {
		base64Image, ok := jsonParsed.Path(path).Data().(string)
		if !ok {
			processingError = errors.New(fmt.Sprintf("Could not find path: %s", path))
			return nil, processingError
		}

		_, err := jsonParsed.Set(nil, path)
		if err != nil {
			processingError = errors.New(fmt.Sprintf("Could not nullify path: %s", path))
			return nil, processingError
		}

		uploadedFile, err := uploadFile(projectId, tempFile{
			path:       path,
			base64File: &base64Image,
		})

		if err != nil {
			processingError = err
			return nil, processingError
		}

		uploadedFiles = append(uploadedFiles, uploadedFile)

		id, err := callback(
			uploadedFile.FileSystemFilePath,
			uploadedFile.Path,
			uploadedFile.MimeType,
			uploadedFile.Extension,
		)

		if err != nil {
			processingError = err
			return nil, processingError
		}

		if err := setJsonFields(jsonParsed, id, uploadedFile); err != nil {
			processingError = err
			return nil, processingError
		}

		uploadedFiles = append(uploadedFiles, uploadedFile)
	}

	return jsonParsed.Bytes(), nil
}
