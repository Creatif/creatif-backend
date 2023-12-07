package queryListByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
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
	c.logBuilder.Add("queryListByID", "Validating...")

	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("queryListByID", "Validated")

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

func (c Main) Logic() (declarations.ListVariable, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("queryListByID", err.Error())
		return declarations.ListVariable{}, appErrors.NewApplicationError(err).AddError("queryListByID.Logic", nil)
	}

	sql := fmt.Sprintf(`
			SELECT lv.id, lv.name, lv.index, lv.behaviour, lv.short_id, lv.metadata, lv.value, lv.groups, lv.created_at, lv.updated_at
			FROM %s AS lv INNER JOIN %s AS l
			ON l.project_id = ? AND (l.name = ? OR l.short_id = ? OR l.id = ?) AND lv.list_id = l.id AND (lv.id = ? OR lv.short_id = ?) AND l.locale_id = ?`,
		(declarations.ListVariable{}).TableName(), (declarations.List{}).TableName())

	var variable declarations.ListVariable
	res := storage.Gorm().
		Raw(sql, c.model.ProjectID, c.model.Name, c.model.Name, c.model.Name, c.model.ItemID, c.model.ItemID, localeID).
		Scan(&variable)

	if res.Error != nil {
		c.logBuilder.Add("queryListByID", res.Error.Error())
		return declarations.ListVariable{}, appErrors.NewDatabaseError(res.Error).AddError("queryListByID.Logic", nil)
	}

	if res.RowsAffected == 0 {
		c.logBuilder.Add("queryListByID", "No rows returned. 404 status code.")
		return declarations.ListVariable{}, appErrors.NewNotFoundError(errors.New("No rows found")).AddError("queryListByID.Logic", nil)
	}

	return variable, nil
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

	return newView(model, c.model.Locale), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, declarations.ListVariable] {
	logBuilder.Add("queryListByID", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
