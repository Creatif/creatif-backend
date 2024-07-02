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

func setJsonFields(jsonParsed *gabs.Container, fileId string, file createdFile) error {
	_, err := jsonParsed.Object(file.Path)
	if err != nil {
		return err
	}

	paths := map[string]string{
		fmt.Sprintf("%s.id", file.Path):        fileId,
		fmt.Sprintf("%s.path", file.Path):      file.PublicFilePath,
		fmt.Sprintf("%s.mimeType", file.Path):  file.MimeType,
		fmt.Sprintf("%s.extension", file.Path): file.Extension,
	}

	for p, v := range paths {
		_, err := jsonParsed.SetP(
			v,
			p,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
