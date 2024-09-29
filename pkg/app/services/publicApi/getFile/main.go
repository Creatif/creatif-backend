package getFile

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/published"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	return nil
}

func (c Main) Authenticate() error {
	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (published.PublishedFile, error) {
	version, err := getVersion(c.model.ProjectID, c.model.Version)
	if err != nil {
		return published.PublishedFile{}, err
	}

	var file published.PublishedFile
	sql := fmt.Sprintf("SELECT id, name, mime_type FROM %s WHERE project_id = ? AND id = ? AND version_id = ?", (published.PublishedFile{}).TableName())

	res := storage.Gorm().Raw(sql, c.model.ProjectID, c.model.FileID, version.ID).Scan(&file)

	if res.Error != nil {
		return published.PublishedFile{}, appErrors.NewApplicationError(res.Error)
	}

	return file, nil
}

func (c Main) Handle() (View, error) {
	if err := c.Validate(); err != nil {
		return View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return View{}, err
	}

	if err := c.Authorize(); err != nil {
		return View{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return View{}, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication) pkg.Job[Model, View, published.PublishedFile] {
	return Main{model: model, auth: auth}
}
