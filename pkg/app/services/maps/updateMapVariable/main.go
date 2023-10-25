package updateMapVariable

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"strings"
)

type Main struct {
	model Model
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

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

func (c Main) Logic() (LogicResult, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		return LogicResult{}, appErrors.NewApplicationError(err).AddError("updateMapVariable.Logic", nil)
	}

	var m declarations.Map
	if res := storage.Gorm().Where("name = ? AND project_id = ? AND locale_id = ?", c.model.MapName, c.model.ProjectID, localeID).First(&m); res.Error != nil {
		return LogicResult{}, appErrors.NewNotFoundError(res.Error).AddError("updateMapVariable.Logic", nil)
	}

	pqGroups := pq.StringArray{}
	if c.model.Values.Groups == nil {
		pqGroups = pq.StringArray{}
	} else {
		for _, v := range c.model.Values.Groups {
			pqGroups = append(pqGroups, v)
		}
	}

	placeholders := make(map[string]interface{})
	updateableFields := ""
	for idx, value := range c.model.Fields {
		var field string
		if value == "name" {
			field = "name = @newName"
			placeholders["newName"] = c.model.Values.Name
		}

		if value == "behaviour" {
			field = "behaviour = @behaviour"
			placeholders["behaviour"] = c.model.Values.Behaviour
		}

		if value == "groups" {
			field = "groups = @groups"
			// HACK: named parameters do not support casting to []text and gorm does not do that but
			// destructures every entry in the array into its parts
			start := "{"
			for i, g := range c.model.Values.Groups {
				start += g
				if i != len(c.model.Values.Groups)-1 {
					start += ","
				}
			}
			start += "}"
			placeholders["groups"] = start
		}

		if value == "metadata" {
			field = "metadata = @metadata"
			placeholders["metadata"] = c.model.Values.Metadata
		}

		if value == "value" {
			field = "value = @value"
			placeholders["value"] = c.model.Values.Value
		}

		updateableFields += field
		if idx != len(c.model.Fields)-1 {
			updateableFields += ","
		}
	}

	fmt.Println(updateableFields)

	placeholders["name"] = c.model.VariableName
	placeholders["mapID"] = m.ID
	placeholders["localeID"] = localeID

	returningFields := []string{"id", "short_id", "name", "behaviour", "metadata", "groups", "value", "created_at", "updated_at"}
	var model declarations.MapVariable
	if res := storage.Gorm().Raw(fmt.Sprintf(
		"UPDATE %s SET %s WHERE name = @name AND map_id = @mapID AND locale_id = @localeID RETURNING %s", (declarations.MapVariable{}).TableName(), updateableFields, strings.Join(returningFields, ",")),
		placeholders,
	).Scan(&model); res.Error != nil || res.RowsAffected == 0 {
		return LogicResult{}, appErrors.NewNotFoundError(errors.New("Could not update map")).AddError("updateMapVariable.Logic", nil)
	}

	return LogicResult{
		Map:    m,
		Locale: c.model.Locale,
		Entry:  model,
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

func New(model Model) pkg.Job[Model, View, LogicResult] {
	return Main{model: model}
}
