package getFile

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
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
	if err := c.auth.Authenticate(); err != nil {
		return appErrors.NewAuthenticationError(err)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (declarations.File, error) {
	var image declarations.File
	sql := fmt.Sprintf("SELECT id, name, field_name, mime_type, extension FROM %s WHERE project_id = ? AND id = ?", (declarations.File{}).TableName())

	if res := storage.Gorm().Raw(sql, c.model.ProjectID, c.model.StructureID).Scan(&image); res.Error != nil {
		return declarations.File{}, appErrors.NewApplicationError(res.Error)
	}

	return image, nil
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, View, declarations.File] {
	return Main{model: model, auth: auth}
}
