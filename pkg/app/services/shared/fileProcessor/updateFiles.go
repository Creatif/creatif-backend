package fileProcessor

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/events"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"os"
)

func UpdateFiles(
	projectId string,
	value []byte,
	imagePaths []string,
	currentImages []declarations.Image,
	createCallback callbackCreateFn,
	updateCallback callbackUpdateFn,
	deleteCallback callbackDeleteFn,
) ([]byte, error) {
	uploadedPaths := make([]string, 0)
	var processingError error
	defer func() {
		if processingError == nil {
			return
		}

		for _, path := range uploadedPaths {
			if err := os.Remove(path); err != nil {
				events.DispatchEvent(events.NewFileNotRemoveEvent(path, "", projectId))
			}
		}
	}()

	for _, currentImage := range currentImages {
		raw := gjson.GetBytes(value, currentImage.FieldName)
		existsInPath := sdk.IncludesFn(imagePaths, func(item string) bool {
			return currentImage.FieldName == item
		})

		// if the file has not been sent in request but exists in db, it means that it is removed. Remove it here
		if raw.Type == gjson.Null {
			if err := os.Remove(currentImage.Name); err != nil {
				events.DispatchEvent(events.NewFileNotRemoveEvent(currentImage.Name, "", projectId))
				return nil, err
			}

			if err := deleteCallback(currentImage.ID); err != nil {
				return nil, err
			}

			continue
		}

		if !existsInPath && raw.Type != gjson.Null {
			newValue, err := sjson.DeleteBytes(value, currentImage.FieldName)
			if err != nil {
				return nil, err
			}
			value = newValue

			if err := os.Remove(currentImage.Name); err != nil {
				events.DispatchEvent(events.NewFileNotRemoveEvent(currentImage.Name, "", projectId))
				return nil, err
			}

			if err := deleteCallback(currentImage.ID); err != nil {
				return nil, err
			}

			continue
		}

		base64Image := gjson.GetBytes(value, currentImage.FieldName).Str
		newValue, err := sjson.SetBytes(value, currentImage.FieldName, nil)
		value = newValue

		uploadedFile, err := uploadFile(projectId, tempFile{
			path:       currentImage.FieldName,
			base64File: &base64Image,
		})

		if err != nil {
			processingError = err
			return nil, processingError
		}

		value, err = setJsonFields(value, currentImage.ID, uploadedFile)

		if err := updateCallback(currentImage.ID, uploadedFile.FileSystemFilePath, uploadedFile.Path, uploadedFile.MimeType, uploadedFile.Extension); err != nil {
			processingError = err
			return nil, processingError
		}

		if err := os.Remove(currentImage.Name); err != nil {
			events.DispatchEvent(events.NewFileNotRemoveEvent(currentImage.Name, "", projectId))
			processingError = err
			return nil, processingError
		}

		uploadedPaths = append(uploadedPaths, uploadedFile.FileSystemFilePath)
	}

	uploadedPaths = make([]string, 0)
	for _, uploadingPath := range imagePaths {
		exists := sdk.IncludesFn(currentImages, func(item declarations.Image) bool {
			return item.FieldName == uploadingPath
		})

		if exists {
			continue
		}

		raw := gjson.GetBytes(value, uploadingPath)
		if raw.Type == gjson.Null {
			processingError = errors.New(fmt.Sprintf("Uploading path %s does not exist", uploadingPath))
			return nil, processingError
		}

		base64Image := raw.Str
		modifiedValue, err := sjson.SetBytes(value, uploadingPath, nil)
		if err != nil {
			processingError = err
			return nil, processingError
		}

		if base64Image == "" {
			continue
		}

		uploadedFile, err := uploadFile(projectId, tempFile{
			path:       uploadingPath,
			base64File: &base64Image,
		})

		if err != nil {
			processingError = err
			return nil, processingError
		}

		id, err := createCallback(
			uploadedFile.FileSystemFilePath,
			uploadedFile.Path,
			uploadedFile.MimeType,
			uploadedFile.Extension,
		)

		if err != nil {
			processingError = err
			return nil, processingError
		}

		newValue, err := setJsonFields(modifiedValue, id, uploadedFile)
		value = newValue

		uploadedPaths = append(uploadedPaths, uploadedFile.FileSystemFilePath)
	}

	return value, nil
}
