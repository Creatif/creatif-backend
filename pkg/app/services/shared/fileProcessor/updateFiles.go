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

/*
*
Updating reuploads all the files again. That means everything that was before gets deleted and
all new is created.
*/
func UpdateFiles(
	projectId string,
	value []byte,
	imagePaths []string,
	currentFiles []declarations.File,
	createCallback callbackCreateFn,
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

	groupedFiles := createFileGrouping(currentFiles)
	uploadedPaths = make([]string, 0)

	for fieldName, files := range groupedFiles {
		shouldSkip, newValue, err := doDeleteIfNotExists(projectId, files, imagePaths, value, deleteCallback)
		if err != nil {
			return nil, err
		}

		/**
		Skipping happens if the user removed the field entirely on the frontend or renames the field
		on the frontend.
		*/
		if shouldSkip {
			value = newValue
			continue
		}

		pathValue := gjson.GetBytes(value, fieldName)

		/**
		if this is true, it means that only one file has been uploaded, therefor,
		we can take that file and delete the rest
		*/
		if pathValue.Type == gjson.String {
			firstFile := files[0]

			/**
			Since only a single file has been uploaded, we can safely delete all
			the db files since they are no longer valid and do not exist by field name
			since field name is unique.
			This might happen if the user has multiple=true on frontend and changed
			it. That means that the previous state had multiple files under the same
			field name. All those files can now be removed
			*/
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
				continue
			}

			id, uploadedFile, err := doCreateUpload(projectId, &pathValue.Str, firstFile.FieldName, createCallback)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			modifiedValue, err := setJsonFields(value, id, uploadedFile)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			uploadedPaths = append(uploadedPaths, uploadedFile.FileSystemFilePath)
			value = modifiedValue
		}

		/**
		If this path is entered, it means that this field name has multiple files associated with it.
		In this case, previous files can be safely deleted and new files are then uploaded
		*/
		if pathValue.Type == gjson.JSON {
			firstFile := files[0]
			// checking that this field is cleared on the user end.
			// this means that user pressed the X button or that no values have been uploaded on update.
			// either way, this code must delete all files and db entries associated with this field name
			if len(pathValue.Array()) == 0 {
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

			paths := make([]map[string]string, 0)
			pathValue.ForEach(func(key, result gjson.Result) bool {
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
					"fileName":  uploadedFile.FileName,
				})

				return true
			})

			modifiedValue, err := sjson.SetBytes(value, firstFile.FieldName, paths)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			value = modifiedValue
		}
	}

	for _, uploadingPath := range imagePaths {
		// this check ensures that this path is already handled
		// by the above code that makes a diff based on db files.
		// no need to process in that case so its ok to skip
		exists := sdk.IncludesFn(currentFiles, func(item declarations.File) bool {
			return item.FieldName == uploadingPath
		})

		if exists {
			continue
		}

		pathValue := gjson.GetBytes(value, uploadingPath)

		if pathValue.Type == gjson.Null {
			continue
		}

		if pathValue.Type == gjson.String && pathValue.Str != "" {
			modifiedValue, err := sjson.SetBytes(value, uploadingPath, nil)
			if err != nil {
				processingError = errors.New(fmt.Sprintf("Could not nullify path: %s", uploadingPath))
				return nil, processingError
			}

			id, uploadedFile, err := doCreateUpload(projectId, &pathValue.Str, uploadingPath, createCallback)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			modifiedValue, err = setJsonFields(value, id, uploadedFile)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			uploadedPaths = append(uploadedPaths, uploadedFile.FileSystemFilePath)
			value = modifiedValue
		}

		if pathValue.Type == gjson.JSON && len(pathValue.Array()) != 0 {
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
					"fileName":  uploadedFile.FileName,
				})

				uploadedPaths = append(uploadedPaths, uploadedFile.FileSystemFilePath)

				return true
			})

			modifiedValue, err := sjson.SetBytes(value, uploadingPath, nil)
			if err != nil {
				processingError = errors.New(fmt.Sprintf("Could not nullify path: %s", uploadingPath))
				return nil, processingError
			}

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

/*
*
Deletes files from db and []byte value if they don't exist in the current request value or they don't exist in path
This happens if the user removes the field entirely on the frontend or renames the filed (name attribute).
*/
func doDeleteIfNotExists(projectId string, files []declarations.File, filePaths []string, value []byte, deleteCallback callbackDeleteFn) (bool, []byte, error) {
	firstFile := files[0]

	raw := gjson.GetBytes(value, firstFile.FieldName)
	existsInPath := sdk.IncludesFn(filePaths, func(item string) bool {
		return firstFile.FieldName == item
	})

	// if the file has not been sent in request but exists in db, it means that it is removed. Remove it here
	if raw.Type == gjson.Null {
		if err := os.Remove(firstFile.Name); err != nil {
			events.DispatchEvent(events.NewFileNotRemoveEvent(firstFile.Name, "", projectId))
		}

		for _, f := range files {
			if err := deleteCallback(f.ID, ""); err != nil {
				return false, nil, err
			}
		}

		return true, value, nil
	}

	if !existsInPath && raw.Type != gjson.Null {
		newValue, err := sjson.DeleteBytes(value, firstFile.FieldName)
		if err != nil {
			return false, nil, err
		}
		value = newValue

		for _, f := range files {
			if err := os.Remove(f.Name); err != nil {
				events.DispatchEvent(events.NewFileNotRemoveEvent(f.Name, "", projectId))
			}

			if err := deleteCallback(f.ID, ""); err != nil {
				return false, nil, err
			}
		}
	}

	return false, value, nil
}

/*
*
Groups file by field name into a map. FieldName is the name gotten from the frontend and is the name
of the frontend component. This method maps all field names to its files because there could be multiple
files in one field name.
*/
func createFileGrouping(dbFiles []declarations.File) map[string][]declarations.File {
	m := make(map[string][]declarations.File)

	for _, f := range dbFiles {
		_, ok := m[f.FieldName]
		if !ok {
			m[f.FieldName] = make([]declarations.File, 0)
		}

		m[f.FieldName] = append(m[f.FieldName], f)
	}

	return m
}
