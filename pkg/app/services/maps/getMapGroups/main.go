package getMapGroups

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"strings"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("getMapGroups", "Validating...")

	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("getMapGroups", "Validated")

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

func (c Main) Logic() ([]string, error) {
	sql := fmt.Sprintf(`
SELECT groups FROM %s AS lv 
    INNER JOIN %s AS l ON l.project_id = ? AND lv.map_id = l.id AND (l.name = ? OR l.id = ? OR l.short_id = ?) AND (lv.id = ? OR lv.short_id = ?)
`, (declarations2.MapVariable{}).TableName(), (declarations2.Map{}).TableName())
	var duplicatedModel []LogicModel
	res := storage.Gorm().Raw(sql, c.model.ProjectID, c.model.Name, c.model.Name, c.model.Name, c.model.ItemID, c.model.ItemID).Scan(&duplicatedModel)

	if res.Error != nil && res.RowsAffected == 0 {
		return nil, appErrors.NewNotFoundError(res.Error)
	} else if res.Error != nil && strings.Contains(res.Error.Error(), "cannot accumulate empty arrays") {
		return make([]string, 0), nil
	} else if res.Error != nil {
		return nil, appErrors.NewApplicationError(res.Error)
	}

	distinctModels := make([]string, 0)
	for _, v := range duplicatedModel {
		groups := v.Groups

		for _, g := range groups {
			if !sdk.Includes(distinctModels, g) {
				distinctModels = append(distinctModels, g)
			}
		}
	}

	return distinctModels, nil
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, []string, []string] {
	logBuilder.Add("getMapGroups", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
