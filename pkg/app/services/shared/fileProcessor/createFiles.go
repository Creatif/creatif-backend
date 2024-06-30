package fileProcessor

import (
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"os"
)

func UploadFiles(projectId string, value []byte, imagePaths []string, callback callbackCreateFn) ([]byte, []createdFile, error) {
	jsonParsed, err := gabs.ParseJSON(value)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Parsing JSON failed: %s", err.Error()))
	}

	files := make([]fileResult, 0)
	for _, path := range imagePaths {
		base64Image, ok := jsonParsed.Path(path).Data().(string)
		if !ok {
			files = append(files, fileResult{
				createdFile: createdFile{},
				error:       errors.New(fmt.Sprintf("Could not find path: %s", path)),
			})

			break
		}

		_, err := jsonParsed.Set(nil, path)
		if err != nil {
			files = append(files, fileResult{
				createdFile: createdFile{},
				error:       errors.New(fmt.Sprintf("Could not nullify path: %s", path)),
			})

			break
		}

		uploadedFile, err := uploadFile(projectId, tempFile{
			path:       path,
			base64File: &base64Image,
		})

		if err != nil {
			files = append(files, fileResult{
				createdFile: createdFile{},
				error:       err,
			})

			break
		}

		id, err := callback(
			uploadedFile.FileSystemFilePath,
			uploadedFile.Path,
			uploadedFile.MimeType,
			uploadedFile.Extension,
		)

		if err != nil {
			files = append(files, fileResult{
				createdFile: createdFile{},
				error:       err,
			})

			break
		}

		uploadedFile.ID = id
		files = append(files, fileResult{
			createdFile: uploadedFile,
			error:       nil,
		})
	}

	for _, uploadedFile := range files {
		if uploadedFile.error != nil {
			return value, nil, uploadedFile.error
		}
	}

	createdFiles := make([]createdFile, 0)
	for _, uploadedFile := range files {
		if uploadedFile.error != nil {
			for _, createdFile := range createdFiles {
				os.Remove(createdFile.FileSystemFilePath)
			}

			return value, createdFiles, uploadedFile.error
		}

		_, err := jsonParsed.Object(uploadedFile.createdFile.Path)
		if err != nil {
			if err != nil {
				for _, createdFile := range createdFiles {
					os.Remove(createdFile.FileSystemFilePath)
				}

				return value, createdFiles, err
			}
		}

		paths := map[string]string{
			fmt.Sprintf("%s.id", uploadedFile.createdFile.Path):        uploadedFile.createdFile.ID,
			fmt.Sprintf("%s.path", uploadedFile.createdFile.Path):      uploadedFile.createdFile.PublicFilePath,
			fmt.Sprintf("%s.mimeType", uploadedFile.createdFile.Path):  uploadedFile.createdFile.MimeType,
			fmt.Sprintf("%s.extension", uploadedFile.createdFile.Path): uploadedFile.createdFile.Extension,
		}

		for p, v := range paths {
			_, err := jsonParsed.SetP(
				v,
				p,
			)

			if err != nil {
				for _, createdFile := range createdFiles {
					os.Remove(createdFile.FileSystemFilePath)
				}

				return value, createdFiles, err
			}
		}

		createdFiles = append(createdFiles, uploadedFile.createdFile)
	}

	return jsonParsed.Bytes(), createdFiles, nil
}
