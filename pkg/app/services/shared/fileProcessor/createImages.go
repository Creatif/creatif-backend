package fileProcessor

import (
	"creatif/pkg/lib/constants"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/google/uuid"
	"os"
	"regexp"
	"strings"
	"sync"
)

type tempFile struct {
	path       string
	base64File *string
}

type createdFile struct {
	ID                 string
	Path               string
	Extension          string
	MimeType           string
	PublicFilePath     string
	FileSystemFilePath string
}

type fileResult struct {
	createdFile createdFile
	error       error
}

type callbackCreateFn = func(fileSystemFilePath, path, mimeType, extension string) (string, error)

func UploadFiles(projectId string, value []byte, imagePaths []string, callback callbackCreateFn) ([]byte, []createdFile, error) {
	jsonParsed, err := gabs.ParseJSON(value)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Parsing JSON failed: %s", err.Error()))
	}

	files := make([]fileResult, 0)
	wg := &sync.WaitGroup{}
	for _, path := range imagePaths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			base64Image, ok := jsonParsed.Path(path).Data().(string)
			if !ok {
				files = append(files, fileResult{
					createdFile: createdFile{},
					error:       errors.New(fmt.Sprintf("Could not find path: %s", path)),
				})

				return
			}

			_, err := jsonParsed.Set(nil, path)
			if err != nil {
				files = append(files, fileResult{
					createdFile: createdFile{},
					error:       errors.New(fmt.Sprintf("Could not nullify path: %s", path)),
				})

				return
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

				return
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

				return
			}

			uploadedFile.ID = id
			files = append(files, fileResult{
				createdFile: uploadedFile,
				error:       nil,
			})
		}(path)
	}

	wg.Wait()

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

func uploadFile(projectId string, file tempFile) (createdFile, error) {
	spl := strings.Split(*file.base64File, "base64,")
	dec, err := base64.StdEncoding.DecodeString(spl[1])
	if err != nil {
		return createdFile{}, fmt.Errorf("Could not decode base64 image: %w", err)
	}

	mimeType, extension, err := extractAndValidateMimeType(file.base64File)
	if err != nil {
		return createdFile{}, err
	}

	fileName := fmt.Sprintf("%s.%s", uuid.NewString(), extension)
	filePath := fmt.Sprintf(
		"%s/%s/%s",
		constants.AssetsDirectory,
		projectId,
		fileName,
	)

	f, err := os.Create(filePath)

	if err != nil {
		return createdFile{}, err
	}

	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return createdFile{}, err
	}

	if err := f.Sync(); err != nil {
		return createdFile{}, err
	}

	return createdFile{
		Path:               file.path,
		MimeType:           mimeType,
		Extension:          extension,
		PublicFilePath:     fmt.Sprintf("/api/v1/static/%s/%s", projectId, fileName),
		FileSystemFilePath: filePath,
	}, nil
}

func extractAndValidateMimeType(image *string) (string, string, error) {
	re := regexp.MustCompile(`data:(.*);`)
	match := re.FindStringSubmatch(*image)
	if match == nil {
		return "", "", errors.New("Could not determine mime type")
	}

	if len(match) < 2 {
		return "", "", errors.New("Could not determine mime type")
	}

	mimeType := match[1]
	sep := strings.Split(mimeType, "#")
	if len(sep) < 2 {
		return "", "", errors.New("base64 has mime type but could not determine extension")
	}

	return sep[0], sep[1], nil
}
