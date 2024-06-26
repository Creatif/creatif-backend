package updateList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("updateList", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}
	c.logBuilder.Add("updateList", "Validated")

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

func (c Main) Logic() (declarations.List, error) {
	var existing declarations.List
	if res := storage.Gorm().Where(fmt.Sprintf("(name = ? OR id = ? OR short_id = ?) AND project_id = ?"), c.model.Name, c.model.Name, c.model.Name, c.model.ProjectID).First(&existing); res.Error != nil {
		c.logBuilder.Add("updateList", res.Error.Error())

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return declarations.List{}, appErrors.NewNotFoundError(res.Error).AddError("updateList.Logic", nil)
		}

		return declarations.List{}, appErrors.NewDatabaseError(res.Error).AddError("updateList.Logic", nil)
	}

	for _, f := range c.model.Fields {
		if f == "name" {
			existing.Name = c.model.Values.Name
		}
	}

	var updated declarations.List
	if res := storage.Gorm().Model(&updated).Clauses(clause.Returning{Columns: []clause.Column{
		{Name: "id"},
		{Name: "project_id"},
		{Name: "name"},
		{Name: "short_id"},
		{Name: "created_at"},
		{Name: "updated_at"},
	}}).Where("id = ?", existing.ID).Select(c.model.Fields).Updates(existing); res.Error != nil {
		c.logBuilder.Add("updateList", res.Error.Error())
		return declarations.List{}, appErrors.NewApplicationError(res.Error).AddError("updateList.Logic", nil)
	}

	return updated, nil
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, declarations.List] {
	logBuilder.Add("updateList", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
