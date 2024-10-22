package createProject

import (
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/constants"
	"fmt"
	"os"
)

func createProjectInPublicDirectory(id string) error {
	projectPublicDir := fmt.Sprintf("%s/%s", constants.PublicDirectory, id)
	_, err := os.Stat(projectPublicDir)
	if !os.IsNotExist(err) {
		return appErrors.NewApplicationError(err).AddError("createProject", nil)
	}

	if os.IsNotExist(err) {
		err := os.MkdirAll(projectPublicDir, os.ModeDir)

		if err != nil {
			return appErrors.NewApplicationError(err).AddError("createProject", nil)
		}
	}

	return nil
}

func createProjectInAssetsDirectory(id string) error {
	projectPublicDir := fmt.Sprintf("%s/%s", constants.AssetsDirectory, id)
	_, err := os.Stat(projectPublicDir)
	if !os.IsNotExist(err) {
		return appErrors.NewApplicationError(err).AddError("createProject", nil)
	}

	if os.IsNotExist(err) {
		err := os.MkdirAll(projectPublicDir, os.ModePerm)

		if err != nil {
			return appErrors.NewApplicationError(err).AddError("createProject", nil)
		}
	}

	return nil
}
