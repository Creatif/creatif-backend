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
			continue
		}

		if pathValue.Type == gjson.String {
			modifiedValue, err := sjson.SetBytes(value, path, nil)
			if err != nil {
				processingError = errors.New(fmt.Sprintf("Could not nullify path: %s", path))
				return nil, processingError
			}

			id, uploadedFile, err := doCreateUpload(projectId, &pathValue.Str, path, callback)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			modifiedValue, err = setJsonFields(value, id, uploadedFile)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			uploadedFiles = append(uploadedFiles, uploadedFile)
			value = modifiedValue

			return value, nil
		}

		if pathValue.Type == gjson.JSON {
			// 1. upload files and create temp json fields
			// 2. set json fields
			paths := make([]map[string]string, 0)
			pathValue.ForEach(func(key, result gjson.Result) bool {
				id, uploadedFile, err := doCreateUpload(projectId, &result.Str, path, callback)
				if err != nil {
					processingError = err
					return false
				}

				paths = append(paths, map[string]string{
					"id":        id,
					"path":      uploadedFile.PublicFilePath,
					"mimeType":  uploadedFile.MimeType,
					"extension": uploadedFile.Extension,
				})

				uploadedFiles = append(uploadedFiles, uploadedFile)

				return true
			})

			modifiedValue, err := sjson.SetBytes(value, path, nil)
			if err != nil {
				processingError = errors.New(fmt.Sprintf("Could not nullify path: %s", path))
				return nil, processingError
			}

			modifiedValue, err = sjson.SetBytes(modifiedValue, path, paths)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			value = modifiedValue
		}
	}

	return value, nil
}

func doCreateUpload(projectId string, uploadingBase64 *string, path string, callback callbackCreateFn) (string, createdFile, error) {
	uploadedFile, err := uploadFile(projectId, tempFile{
		path:       path,
		base64File: uploadingBase64,
	})

	if err != nil {
		return "", createdFile{}, err
	}

	id, err := callback(
		uploadedFile.FileSystemFilePath,
		uploadedFile.Path,
		uploadedFile.MimeType,
		uploadedFile.Extension,
	)

	if err != nil {
		return "", createdFile{}, err
	}

	return id, uploadedFile, nil
}
