package fileProcessor

import (
	"creatif/pkg/app/services/events"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"os"
)

func UploadFiles(projectId string, value []byte, imagePaths []string, callback callbackCreateFn) ([]byte, error) {
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
		pathValue := gjson.GetBytes(value, path)

		if pathValue.Type == gjson.Null {
			processingError = errors.New(fmt.Sprintf("Could not find path: %s", path))
			return nil, processingError
		}

		modifiedValue, err := sjson.SetBytes(value, path, nil)
		if err != nil {
			processingError = errors.New(fmt.Sprintf("Could not nullify path: %s", path))
			return nil, processingError
		}

		uploadedFile, err := uploadFile(projectId, tempFile{
			path:       path,
			base64File: &pathValue.Str,
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

		modifiedValue, err = setJsonFields(modifiedValue, id, uploadedFile)
		if err != nil {
			processingError = err
			return nil, processingError
		}

		uploadedFiles = append(uploadedFiles, uploadedFile)
		value = modifiedValue
	}

	return value, nil
}
