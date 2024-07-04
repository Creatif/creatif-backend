package fileProcessor

import (
	"creatif/pkg/lib/constants"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/tidwall/sjson"
	"os"
	"regexp"
	"strings"
)

func uploadFile(projectId string, file tempFile) (createdFile, error) {
	spl := strings.Split(*file.base64File, "base64,")
	if len(spl) != 2 {
		return createdFile{}, errors.New(fmt.Sprintf("Path %s invalid. Invalid base64. No data path", file.path))
	}
	dec, err := base64.StdEncoding.DecodeString(spl[1])
	if err != nil {
		return createdFile{}, err
	}

	fmt.Println("extracting mime type and extension")
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

	fmt.Println("creating file on filesystem")
	f, err := os.Create(filePath)

	if err != nil {
		return createdFile{}, err
	}

	defer f.Close()

	fmt.Println("writing file")
	if _, err := f.Write(dec); err != nil {
		return createdFile{}, err
	}

	fmt.Println("syncing file")
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

func setJsonFields(value []byte, fileId string, file createdFile) ([]byte, error) {
	paths := map[string]string{
		"id":        fileId,
		"path":      file.PublicFilePath,
		"mimeType":  file.MimeType,
		"extension": file.Extension,
	}

	return sjson.SetBytes(value, file.Path, paths)
}

func replacePath(path string) string {
	if strings.Contains(path, ".") {
		path = strings.Replace(path, ".", "/", -1)
		return fmt.Sprintf("/%s", path)
	}

	return fmt.Sprintf("/%s", path)
}
