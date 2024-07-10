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

	// questions:
	// 1. What if a field was a single file upload but changes to multiple files?
	// 2. What if a field was a multiple file upload but becomes a single file upload?
	// 3. What if multiple files changes with more or less files uploaded?

	// answers:
	// 1. That single file must be deleted along with db entry and multiple files must be created
	// 2. Those multiple files must be deleted and a single entry must be created
	// 3. Files that exist, must be updated. Files that don't exist anymore must be deleted. Files that are new, must be created

	// 1. Create a map with field name -> file list grouping

	groupedFiles := createFileGrouping(currentImages)
	uploadedPaths = make([]string, 0)

	fmt.Println(groupedFiles)

	for fieldName, files := range groupedFiles {
		fmt.Println("processing field name: ", fieldName, len(files))
		firstFile := files[0]
		/*		shouldSkip, newValue, err := doDeleteIfNotExists(projectId, firstFile.ID, imagePaths, value, firstFile.Name, firstFile.FieldName, deleteCallback)
				if err != nil {
					return nil, err
				}

				if shouldSkip {
					value = newValue
					continue
				}*/

		pathValue := gjson.GetBytes(value, fieldName)

		if pathValue.Type == gjson.String {
			fmt.Println("path value string", fieldName)
			for _, f := range files {
				if err := deleteCallback(f.ID, f.FieldName); err != nil {
					processingError = err
					return nil, processingError
				}

				if err := os.Remove(f.Name); err != nil {
					events.DispatchEvent(events.NewFileNotRemoveEvent(f.Name, "", projectId))
				}
			}

			if pathValue.Str == "" {
				fmt.Println("path value empty, continue: ", fieldName)
				continue
			}

			id, uploadedFile, err := doCreateUpload(projectId, &pathValue.Str, firstFile.FieldName, createCallback)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			fmt.Println("upload created for string", fieldName)
			modifiedValue, err := setJsonFields(value, id, uploadedFile)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			uploadedPaths = append(uploadedPaths, uploadedFile.FileSystemFilePath)
			value = modifiedValue
		}

		if pathValue.Type == gjson.JSON {
			fmt.Println("path value array", fieldName)
			// checking that this field is cleared on the user end.
			// this means that user pressed the X button or that no values have been uploaded on update
			// either way, this code must delete all images and db entries associated with this field name
			if len(pathValue.Array()) == 0 {
				fmt.Println("array path value empty", fieldName)
				for _, f := range files {
					if err := deleteCallback(f.ID, f.FieldName); err != nil {
						processingError = err
						return nil, processingError
					}

					if err := os.Remove(f.Name); err != nil {
						events.DispatchEvent(events.NewFileNotRemoveEvent(f.Name, "", projectId))
					}
				}

				continue
			}

			for _, f := range files {
				if err := deleteCallback(f.ID, f.FieldName); err != nil {
					processingError = err
					return nil, processingError
				}

				if err := os.Remove(f.Name); err != nil {
					events.DispatchEvent(events.NewFileNotRemoveEvent(f.Name, "", projectId))
				}
			}

			fmt.Println("created array paths", fieldName)
			paths := make([]map[string]string, 0)
			pathValue.ForEach(func(key, result gjson.Result) bool {
				fmt.Println(key, "creating upload", fieldName)
				id, uploadedFile, err := doCreateUpload(projectId, &result.Str, firstFile.FieldName, createCallback)
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

				return true
			})

			fmt.Println("creating paths", fieldName)
			modifiedValue, err := sjson.SetBytes(value, firstFile.FieldName, paths)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			value = modifiedValue
		}
	}

	for _, uploadingPath := range imagePaths {
		fmt.Println("upload path: ", uploadingPath)
		// this check ensures that this path is already handled
		// by the above code that makes a diff based on db files.
		// no need to process in that case so its ok to skip
		exists := sdk.IncludesFn(currentImages, func(item declarations.Image) bool {
			return item.FieldName == uploadingPath
		})

		if exists {
			fmt.Println("path already processed", uploadingPath)
			continue
		}

		pathValue := gjson.GetBytes(value, uploadingPath)

		if pathValue.Type == gjson.Null {
			fmt.Println("path value is null", uploadingPath)
			continue
		}

		if pathValue.Type == gjson.String && pathValue.Str != "" {
			fmt.Println("path value is string", uploadingPath)
			modifiedValue, err := sjson.SetBytes(value, uploadingPath, nil)
			if err != nil {
				processingError = errors.New(fmt.Sprintf("Could not nullify path: %s", uploadingPath))
				return nil, processingError
			}

			fmt.Println("creating upload", uploadingPath)
			id, uploadedFile, err := doCreateUpload(projectId, &pathValue.Str, uploadingPath, createCallback)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			fmt.Println("setting json fields", uploadingPath)
			modifiedValue, err = setJsonFields(value, id, uploadedFile)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			uploadedPaths = append(uploadedPaths, uploadedFile.FileSystemFilePath)
			value = modifiedValue

			return value, nil
		}

		if pathValue.Type == gjson.JSON && len(pathValue.Array()) != 0 {
			fmt.Println("uploading path is json", uploadingPath)
			// 1. upload files and create temp json fields
			// 2. set json fields
			paths := make([]map[string]string, 0)
			pathValue.ForEach(func(key, result gjson.Result) bool {
				id, uploadedFile, err := doCreateUpload(projectId, &result.Str, uploadingPath, createCallback)
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

				uploadedPaths = append(uploadedPaths, uploadedFile.FileSystemFilePath)

				return true
			})

			fmt.Println("nullify bytes", uploadingPath)
			modifiedValue, err := sjson.SetBytes(value, uploadingPath, nil)
			if err != nil {
				processingError = errors.New(fmt.Sprintf("Could not nullify path: %s", uploadingPath))
				return nil, processingError
			}

			fmt.Println("setting bytes", uploadingPath)
			modifiedValue, err = sjson.SetBytes(modifiedValue, uploadingPath, paths)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			value = modifiedValue
		}
	}

	return value, nil
}

func doUpdateSingle(projectId, fileId, fieldName, filePath string, value []byte, updateCallback callbackUpdateFn) ([]byte, string, error) {
	base64Image := gjson.GetBytes(value, fieldName).Str
	newValue, err := sjson.SetBytes(value, fieldName, nil)
	if err != nil {
		return nil, "", err
	}
	value = newValue

	uploadedFile, err := uploadFile(projectId, tempFile{
		path:       fieldName,
		base64File: &base64Image,
	})

	if err != nil {
		return nil, "", err
	}

	newValue, err = setJsonFields(value, fileId, uploadedFile)
	value = newValue

	if err := updateCallback(fileId, uploadedFile.FileSystemFilePath, uploadedFile.Path, uploadedFile.MimeType, uploadedFile.Extension); err != nil {
		return nil, "", err
	}

	if err := os.Remove(filePath); err != nil {
		events.DispatchEvent(events.NewFileNotRemoveEvent(filePath, "", projectId))
	}

	return value, uploadedFile.FileSystemFilePath, nil
}

/*
*
Deletes files from db and []byte value if they don't exist in the current request value or they don't exist in path
This happens if the user removes the field entirely on the frontend or renames the filed (name attribute).
*/
func doDeleteIfNotExists(projectId, fileId string, filePaths []string, value []byte, filePath, fieldName string, deleteCallback callbackDeleteFn) (bool, []byte, error) {
	raw := gjson.GetBytes(value, fieldName)
	existsInPath := sdk.IncludesFn(filePaths, func(item string) bool {
		return fieldName == item
	})

	// if the file has not been sent in request but exists in db, it means that it is removed. Remove it here
	if raw.Type == gjson.Null {
		if err := os.Remove(filePath); err != nil {
			events.DispatchEvent(events.NewFileNotRemoveEvent(filePath, "", projectId))
			return true, nil, nil
		}

		if err := deleteCallback(fileId, fieldName); err != nil {
			return false, nil, err
		}

		return true, nil, nil
	}

	if !existsInPath && raw.Type != gjson.Null {
		newValue, err := sjson.DeleteBytes(value, fieldName)
		if err != nil {
			return false, nil, err
		}
		value = newValue

		if err := os.Remove(filePath); err != nil {
			events.DispatchEvent(events.NewFileNotRemoveEvent(filePath, "", projectId))
			return true, nil, nil
		}

		if err := deleteCallback(fileId, fieldName); err != nil {
			return false, nil, err
		}
	}

	return true, value, nil
}

func createFileGrouping(dbFiles []declarations.Image) map[string][]declarations.Image {
	m := make(map[string][]declarations.Image)

	for _, f := range dbFiles {
		_, ok := m[f.FieldName]
		if !ok {
			m[f.FieldName] = make([]declarations.Image, 0)
		}

		m[f.FieldName] = append(m[f.FieldName], f)
	}

	return m
}
