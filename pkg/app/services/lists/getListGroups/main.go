package getListGroups

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"fmt"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("queryListByID", "Validating...")

	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("queryListByID", "Validated")

	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return err
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() ([]LogicModel, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("appendToList", err.Error())
		return nil, appErrors.NewApplicationError(err).AddError("getListGroups.Logic", nil)
	}

	varId, varVal := shared.DetermineID("l", c.model.Name, c.model.ID, c.model.ShortID)

	sql := fmt.Sprintf(`
SELECT DISTINCT unnest(ARRAY(
	SELECT groups FROM %s AS lv 
    INNER JOIN %s AS l ON l.project_id = ? AND lv.list_id = l.id AND l.locale_id = ? AND %s
	)
) AS group
`, (declarations2.ListVariable{}).TableName(), (declarations2.List{}).TableName(), varId)
	var model []LogicModel
	res := storage.Gorm().Raw(sql, c.model.ProjectID, localeID, varVal).Scan(&model)

	if res.Error != nil {
		return nil, appErrors.NewNotFoundError(res.Error)
	}

	return model, nil
}

func (c Main) Handle() ([]string, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	if err := c.Authorize(); err != nil {
		return nil, err
	}

	model, err := c.Logic()

	if err != nil {
		return nil, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, []string, []LogicModel] {
	logBuilder.Add("queryListByID", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
