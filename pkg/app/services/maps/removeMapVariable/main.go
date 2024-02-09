package removeMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/shared"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return appErrors.NewAuthenticationError(err)
	}

	err := shared.IsParent(c.model.VariableName)
	if err != nil {
		return appErrors.NewValidationError(map[string]string{
			"isParent": "This variable is a parent and cannot be deleted",
		})
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (interface{}, error) {
	if err := storage.Transaction(func(tx *gorm.DB) error {
		sql := fmt.Sprintf(
			`DELETE FROM %s AS mv USING %s AS m WHERE m.project_id = ? AND mv.map_id = m.id AND (mv.id = ? OR mv.short_id = ?) AND (m.id = ? OR m.short_id = ?)`,
			(declarations.MapVariable{}).TableName(),
			(declarations.Map{}).TableName(),
		)

		res := storage.Gorm().Exec(sql, c.model.ProjectID, c.model.VariableName, c.model.VariableName, c.model.Name, c.model.Name)
		if res.Error != nil {
			c.logBuilder.Add("removeMapVariable", res.Error.Error())
			return res.Error
		}

		if res.RowsAffected == 0 {
			c.logBuilder.Add("removeMapVariable", "No rows returned. Returning 404 status.")
			return res.Error
		}

		if err := shared.RemoveAsParent(c.model.VariableName); err != nil {
			return err
		}
		if err := shared.RemoveAsChild(c.model.VariableName); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, appErrors.NewNotFoundError(err).AddError("removeMapVariable.Logic", nil)
	}

	return nil, nil
}

func (c Main) Handle() (interface{}, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	if err := c.Authorize(); err != nil {
		return nil, err
	}

	_, err := c.Logic()

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, interface{}, interface{}] {
	logBuilder.Add("removeMapVariable", "Created.")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
