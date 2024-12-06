package fileProcessor

import (
	"creatif/pkg/app/services/events"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"os"
)

/*
* UploadFiles

		This function has 3 modes:
	    1. null mode
			This is a mode where the user has a upload field on the frontend but has not uploaded anything to it
			This mode does nothing to modify the underlying json value
	    2. Single string mode
			This means that the user uploaded a single image on the frontend. This mode is entered if the frontend
			does not have 'multiple' option selected, that is the frontend is not uploading multiple files to the backend
		3. JSON mode
			This mode is entered when the user on the frontend has specified 'multiple', that is the user uploads
			multiple files. It is called JSON mode because gjson interprets an array of base64 images as a JSON constant.
			Therefor, this mode iterates through all the files and uploads them one by one.
*/
func UploadFiles(projectId string, value []byte, imagePaths []string, callback callbackCreateFn) ([]byte, error) {
	/**
	uploadedFiles variable exists only if one of the uploaded files has an error but the
	is already uploaded to the filesystem. It works together with the processingError.
	If there is a processingError, the already uploaded file is sent as an event
	to be deleted later. This is done in the events system that is running in intervals.
	At the time of writing this comment, the event system is ran every hour.

	For example, if this function has 5 files to upload, and the third fails, the other two
	will not be processed but the two that have already been processed and written to the
	filesystem and scheduled for removal.

	Note that the file info also written to the database and those but this runs in a transaction
	so if one of the files fails, the database will be safe from invalid entries. This is only
	here if the files are already written to the filesystem. It has nothing to do with the failed
	database entries. This is handled in an atomic transaction.
	*/
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
			/**
			Since we already have the value of the path from gjson.GetBytes(), the value that
			the path is pointing to is nullified to conserve memory since it doesn't have to be
			in memory.
			*/
			modifiedValue, err := sjson.SetBytes(value, path, nil)
			if err != nil {
				processingError = errors.New(fmt.Sprintf("Could not nullify path: %s", path))
				return nil, processingError
			}

			/**
			Since we have the base64 value of the image stored in pathValue variable,
			we can safely create the uploaded file on the filesystem
			*/
			id, uploadedFile, err := doCreateUpload(projectId, &pathValue.Str, path, callback)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			/**
			sets the json fields with the information of the uploaded file. uploadedFile is of type
			createdFile struct and carries all the information of the uploaded file, such as the path
			where it is stored on the files system. See the createdFile struct for more information.
			*/
			modifiedValue, err = setJsonFields(value, id, uploadedFile)
			if err != nil {
				processingError = err
				return nil, processingError
			}

			uploadedFiles = append(uploadedFiles, uploadedFile)
			value = modifiedValue
		}

		if pathValue.Type == gjson.JSON {
			/**
			When pathValue.Type == gjson.JSON, that means that multiple files
			have been uploaded and this part of the code iterates trough all of them
			and does the same this as if it is a single file.
			*/
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
					"fileName":  uploadedFile.FileName,
				})

				uploadedFiles = append(uploadedFiles, uploadedFile)

				return true
			})

			/**
			In this case, the path value is an array of base64 files. We can
			nullify them to conserve memory. This path will then be written
			with all the uploaded paths.
			*/
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
		uploadedFile.FileName,
	)

	if err != nil {
		return "", createdFile{}, err
	}

	return id, uploadedFile, nil
}
