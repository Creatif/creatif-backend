package fileProcessor

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/Jeffail/gabs"
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
	jsonParsed, err := gabs.ParseJSON(value)
	if err != nil {
		fmt.Println("json parse error: ", err)
		return nil, err
	}

	fmt.Println("start upload, current images")
	uploadedPaths := make([]string, 0)
	for _, currentImage := range currentImages {
		existsInJson := jsonParsed.Exists(currentImage.FieldName)
		existsInPath := sdk.IncludesFn(imagePaths, func(item string) bool {
			return currentImage.FieldName == item
		})

		fmt.Println("check exists in json", existsInPath, existsInJson, currentImage.FieldName)

		// if the file has not been sent in request but exists in db, it means that it is removed. Remove it here
		if !existsInJson {
			if err := os.Remove(currentImage.Name); err != nil {
				fmt.Println("os remove error: ", err)
				return nil, err
			}

			if err := deleteCallback(currentImage.ID); err != nil {
				fmt.Println("deleter error: ", err)
				return nil, err
			}

			continue
		}

		if !existsInPath && existsInJson {
			if err := os.Remove(currentImage.Name); err != nil {
				fmt.Println("os remove error: ", err)
				return nil, err
			}

			if err := deleteCallback(currentImage.ID); err != nil {
				fmt.Println("deleter error: ", err)
				return nil, err
			}

			continue
		}

		fmt.Println("get base64")
		base64Image, _ := jsonParsed.Path(currentImage.FieldName).Data().(string)

		fmt.Println("nullify json")
		_, err := jsonParsed.Set(nil, currentImage.FieldName)
		if err != nil {
			fmt.Println("parse set error: ", err)
			return nil, err
		}

		fmt.Println("upload file")
		uploadedFile, err := uploadFile(projectId, tempFile{
			path:       currentImage.FieldName,
			base64File: &base64Image,
		})

		if err != nil {
			fmt.Println("Upload file error: ", err)
			return nil, err
		}

		fmt.Println("create object in json")
		_, err = jsonParsed.Object(currentImage.FieldName)
		if err != nil {
			fmt.Println("parse object error: ", err)
			for _, path := range uploadedPaths {
				os.Remove(path)
			}

			return nil, err
		}

		paths := map[string]string{
			fmt.Sprintf("%s.id", currentImage.FieldName):        currentImage.ID,
			fmt.Sprintf("%s.path", currentImage.FieldName):      uploadedFile.PublicFilePath,
			fmt.Sprintf("%s.mimeType", currentImage.FieldName):  uploadedFile.MimeType,
			fmt.Sprintf("%s.extension", currentImage.FieldName): uploadedFile.Extension,
		}

		fmt.Println("set object paths")
		for p, v := range paths {
			_, err := jsonParsed.SetP(
				v,
				p,
			)

			if err != nil {
				fmt.Println("set p error: ", err)
				for _, path := range uploadedPaths {
					os.Remove(path)
				}

				return nil, err
			}
		}

		fmt.Println("update the db")
		if err := updateCallback(currentImage.ID, uploadedFile.FileSystemFilePath, uploadedFile.Path, uploadedFile.MimeType, uploadedFile.Extension); err != nil {
			fmt.Println("Update error: ", err)
			for _, path := range uploadedPaths {
				os.Remove(path)
			}

			return nil, err
		}

		os.Remove(currentImage.Name)
		uploadedPaths = append(uploadedPaths, uploadedFile.FileSystemFilePath)
	}

	fmt.Println("Start create: ", imagePaths)

	uploadedPaths = make([]string, 0)
	for _, uploadingPath := range imagePaths {
		exists := sdk.IncludesFn(currentImages, func(item declarations.Image) bool {
			return item.FieldName == uploadingPath
		})

		if exists {
			continue
		}

		base64Image, _ := jsonParsed.Path(uploadingPath).Data().(string)
		_, err := jsonParsed.Set(nil, uploadingPath)
		if err != nil {
			return nil, err
		}

		if base64Image == "" {
			continue
		}

		_, err = jsonParsed.Set(nil, uploadingPath)
		if err != nil {
			return nil, err
		}

		uploadedFile, err := uploadFile(projectId, tempFile{
			path:       uploadingPath,
			base64File: &base64Image,
		})

		if err != nil {
			return nil, err
		}

		_, err = jsonParsed.Object(uploadingPath)
		if err != nil {
			for _, path := range uploadedPaths {
				os.Remove(path)
			}

			return nil, err
		}

		id, err := createCallback(
			uploadedFile.FileSystemFilePath,
			uploadedFile.Path,
			uploadedFile.MimeType,
			uploadedFile.Extension,
		)
		if err != nil {
			return nil, err
		}

		paths := map[string]string{
			fmt.Sprintf("%s.id", uploadedFile.Path):        id,
			fmt.Sprintf("%s.path", uploadedFile.Path):      uploadedFile.PublicFilePath,
			fmt.Sprintf("%s.mimeType", uploadedFile.Path):  uploadedFile.MimeType,
			fmt.Sprintf("%s.extension", uploadedFile.Path): uploadedFile.Extension,
		}

		for p, v := range paths {
			_, err := jsonParsed.SetP(
				v,
				p,
			)

			if err != nil {
				for _, path := range uploadedPaths {
					os.Remove(path)
				}

				return nil, err
			}
		}

		uploadedPaths = append(uploadedPaths, uploadedFile.FileSystemFilePath)
	}

	return jsonParsed.Bytes(), nil
}
