package removeMapEntry

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
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
	// user check by project id should be gotten here, with authentication cookie
	var project app.Project
	if err := storage.Get((app.Project{}).TableName(), c.model.ProjectID, &project); err != nil {
		return appErrors.NewAuthenticationError(err).AddError("createVariable.Authenticate", nil)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (interface{}, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("removeMapEntry", err.Error())
		return nil, appErrors.NewApplicationError(err).AddError("removeMapEntry.Logic", nil)
	}

	mapId, mapVal := shared.DetermineID("m", c.model.Name, c.model.MapID, c.model.MapShortID)
	varId, varVal := shared.DetermineID("mv", c.model.VariableName, c.model.VariableID, c.model.VariableShortID)

	sql := fmt.Sprintf(
		`DELETE FROM %s AS mv USING %s AS m WHERE m.project_id = ? AND m.locale_id = ? AND mv.map_id = m.id AND %s AND %s`,
		(declarations.MapVariable{}).TableName(),
		(declarations.Map{}).TableName(),
		mapId,
		varId,
	)

	res := storage.Gorm().Exec(sql, c.model.ProjectID, localeID, mapVal, varVal)
	if res.Error != nil {
		c.logBuilder.Add("removeMapEntry", res.Error.Error())
		return nil, appErrors.NewNotFoundError(res.Error).AddError("removeMapEntry.Logic", nil)
	}

	if res.RowsAffected == 0 {
		c.logBuilder.Add("removeMapEntry", "No rows returned. Returning 404 status.")
		return nil, appErrors.NewNotFoundError(errors.New(fmt.Sprintf("Variable with name '%s' not found.", c.model.VariableName))).AddError("removeMapEntry.Logic", nil)
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
	logBuilder.Add("removeMapEntry", "Created.")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
