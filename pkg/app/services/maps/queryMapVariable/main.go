package queryMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
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
	c.logBuilder.Add("queryMapVariable", "Validating...")

	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("queryMapVariable", "Validated")

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

func (c Main) Logic() (LogicModel, error) {
	sql := fmt.Sprintf(`
			SELECT lv.id, lv.name, lv.behaviour, lv.short_id, lv.metadata, lv.value, lv.groups, lv.created_at, lv.updated_at, lv.locale_id
			FROM %s AS lv INNER JOIN %s AS l
			ON l.project_id = ? AND (l.short_id = ? OR l.id = ?) AND lv.map_id = l.id AND (lv.id = ? OR lv.short_id = ?)`,
		(declarations.MapVariable{}).TableName(), (declarations.Map{}).TableName())

	var variable declarations.MapVariable
	res := storage.Gorm().
		Raw(sql, c.model.ProjectID, c.model.Name, c.model.Name, c.model.ItemID, c.model.ItemID).
		Scan(&variable)

	if res.Error != nil {
		c.logBuilder.Add("queryMapVariable", res.Error.Error())
		return LogicModel{}, appErrors.NewDatabaseError(res.Error).AddError("queryMapVariable.Logic", nil)
	}

	if res.RowsAffected == 0 {
		c.logBuilder.Add("queryMapVariable", "No rows returned. 404 status code.")
		return LogicModel{}, appErrors.NewNotFoundError(errors.New("No rows found")).AddError("queryMapVariable.Logic", nil)
	}

	references, err := shared.QueryReferences(variable.ID, c.model.ProjectID)
	if err != nil {
		return LogicModel{}, err
	}

	for _, r := range references {
		fmt.Println(r.ChildStructureID, "hild")
	}

	return LogicModel{
		Variable:  variable,
		Reference: references,
	}, nil
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, LogicModel] {
	logBuilder.Add("queryMapVariable", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
