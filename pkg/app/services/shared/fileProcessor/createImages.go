package fileProcessor

import (
	"creatif/pkg/app/domain/declarations"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/google/uuid"
	"os"
	"regexp"
	"strings"
)

type tempFile struct {
	path       string
	base64File *string
}

type createdFile struct {
	Path               string
	PublicFilePath     string
	FileSystemFilePath string
}

type updatedFile struct {
	toUpdate []createdFile
	toDelete []createdFile
	toCreate []createdFile
}

func UploadFiles(projectId string, value []byte, imagePaths []string) ([]byte, []createdFile, error) {
	jsonParsed, err := gabs.ParseJSON(value)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Parsing JSON failed: %s", err.Error()))
	}

	base64Files := make([]tempFile, 0)
	for _, path := range imagePaths {
		base64Image, ok := jsonParsed.Path(path).Data().(string)
		if !ok {
			return nil, nil, errors.New(fmt.Sprintf("Could not find path: %s", path))
		}

		if base64Image != "" {
			base64Files = append(base64Files, tempFile{
				path:       path,
				base64File: &base64Image,
			})
		}

		_, err := jsonParsed.Set(nil, path)
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("Could not nullify path: %s", path))
		}
	}

	createdFiles, err := uploadFiles(projectId, base64Files)
	if err != nil {
		for _, file := range createdFiles {
			os.Remove(file.FileSystemFilePath)
		}

		return nil, nil, err
	}

	for _, file := range createdFiles {
		_, err := jsonParsed.Set(file.PublicFilePath, file.Path)
		if err != nil {
			for _, file := range createdFiles {
				os.Remove(file.FileSystemFilePath)
			}

			return value, nil, err
		}
	}

	return jsonParsed.Bytes(), createdFiles, nil
}

func UpdateFiles(projectId string, value []byte, updatedPaths []string, deletedPaths []string, dbImages []declarations.Image) ([]byte, []updatedFile, error) {
	base64Files := make([]tempFile, 0)
	for _, path := range updatedPaths {
		base64Image, ok := jsonParsed.Path(path).Data().(string)
		if !ok {
			return nil, nil, errors.New(fmt.Sprintf("Could not find path: %s", path))
		}

		if base64Image != "" {
			base64Files = append(base64Files, tempFile{
				path:       path,
				base64File: &base64Image,
			})
		}

		_, err := jsonParsed.Set(nil, path)
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("Could not nullify path: %s", path))
		}
	}
}

func uploadFiles(projectId string, base64Files []tempFile) ([]createdFile, error) {
	createdFiles := make([]createdFile, 0)

	for _, file := range base64Files {
		spl := strings.Split(*file.base64File, "base64,")
		dec, err := base64.StdEncoding.DecodeString(spl[1])
		if err != nil {
			return createdFiles, fmt.Errorf("Could not decode base64 image: %w", err)
		}

		extension, err := extractAndValidateMimeType(file.base64File)
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
			Path:               file.path,
			PublicFilePath:     fmt.Sprintf("/api/v1/static/%s/%s", projectId, fileName),
			FileSystemFilePath: filePath,
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
