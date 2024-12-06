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

/*
*
uploadFile is processing the base64 file gotten from the frontend.
base64 uploaded string is uploaded in the following format

data:image/webp#webp;base64,<base64 uploaded file>

There for, the mime type of the uploaded image is in this format and is extracted from it.

 1. Check if the base64 is valid base64 that the backend expects from the frontend
 2. Decode the base64 of the file
 3. Extract and validate the mime type of time file
 4. fileName and filePath is created with the uuid package and the extracted extension
 5. The file is then written to the filesystem and the createdFile struct is then returned
    with all the data that the client caller of this function needs
*/
func uploadFile(projectId string, file tempFile) (createdFile, error) {
	// check if it is a valid base64 from the frontend
	spl := strings.Split(*file.base64File, "base64,")
	if len(spl) != 2 {
		return createdFile{}, errors.New(fmt.Sprintf("Path %s invalid. Invalid base64. No data path", file.path))
	}

	// decode the pure base64 of the image
	dec, err := base64.StdEncoding.DecodeString(spl[1])
	if err != nil {
		return createdFile{}, err
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
		FileName:           fileName,
	}, nil
}

/*
*
the format for extraction is as follows:

data:image/webp#webp;

In case of this example, the returned values are 'image/webp, webp, nil'
where the mimeType is image/webp and extension is webp
*/
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

/*
*
This function is called after the file has already been created and uploaded to the filesystem.
The createdFile struct is then used to replace the image path provided by the frontend
with all the data needed to show the file on the frontend. Note that the underlying json
value is modified here and the results of this modification are returned as bytes
*/
func setJsonFields(value []byte, fileId string, file createdFile) ([]byte, error) {
	paths := map[string]string{
		"id":        fileId,
		"path":      file.PublicFilePath,
		"mimeType":  file.MimeType,
		"extension": file.Extension,
		"fileName":  file.FileName,
	}

	return sjson.SetBytes(value, file.Path, paths)
}
