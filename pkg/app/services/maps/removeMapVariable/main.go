package removeMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
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

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (interface{}, error) {
	sql := fmt.Sprintf(
		`DELETE FROM %s AS mv USING %s AS m WHERE m.project_id = ? AND mv.map_id = m.id AND (mv.id = ? OR mv.short_id = ?) AND (m.id = ? OR m.short_id = ? OR m.name = ?)`,
		(declarations.MapVariable{}).TableName(),
		(declarations.Map{}).TableName(),
	)

	res := storage.Gorm().Exec(sql, c.model.ProjectID, c.model.VariableName, c.model.VariableName, c.model.Name, c.model.Name, c.model.Name)
	if res.Error != nil {
		c.logBuilder.Add("removeMapVariable", res.Error.Error())
		return nil, appErrors.NewNotFoundError(res.Error).AddError("removeMapVariable.Logic", nil)
	}

	if res.RowsAffected == 0 {
		c.logBuilder.Add("removeMapVariable", "No rows returned. Returning 404 status.")
		return nil, appErrors.NewNotFoundError(errors.New(fmt.Sprintf("Variable with name '%s' not found.", c.model.VariableName))).AddError("removeMapVariable.Logic", nil)
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