package addGroups

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
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
	c.logBuilder.Add("addToList", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}
	c.logBuilder.Add("addToList", "Validated.")

	toCreateGroups := make([]string, 0)
	for _, g := range c.model.Groups {
		if g.Action == "create" && g.ID == "" {
			toCreateGroups = append(toCreateGroups, g.Name)
		}
	}

	groupExists := fmt.Sprintf("SELECT id FROM %s WHERE project_id = ? AND name IN(?)", (declarations.Group{}).TableName())

	var groups []string
	if res := storage.Gorm().Raw(groupExists, c.model.ProjectID, toCreateGroups).Scan(&groups); res.Error != nil {
		return appErrors.NewApplicationError(res.Error)
	}

	if len(groups) > 0 {
		return appErrors.NewValidationError(map[string]string{
			"groupExists": "Some of the groups exist.",
		})
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

func (c Main) Logic() ([]declarations.Group, error) {
	if transactionErr := storage.Transaction(func(tx *gorm.DB) error {
		if len(c.model.Groups) == 0 {
			res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE project_id = ?", (declarations.Group{}).TableName()), c.model.ProjectID)
			if res.Error != nil {
				return res.Error
			}

			return nil
		}

		var existingGroups []declarations.Group
		res := tx.Raw(fmt.Sprintf("SELECT name, id FROM %s WHERE project_id = ?", (declarations.Group{}).TableName()), c.model.ProjectID).Scan(&existingGroups)
		if res.Error != nil {
			return res.Error
		}

		toCreateGroups := make([]declarations.Group, 0)
		toDeleteGroups := make([]string, 0)
		for _, g := range c.model.Groups {
			if g.Action == "create" && g.ID == "" {
				toCreateGroups = append(toCreateGroups, declarations.NewGroup(c.model.ProjectID, g.Name))
			}

			if g.Action == "remove" && g.ID != "" && sdk.IncludesFn(existingGroups, func(item declarations.Group) bool {
				return item.ID == g.ID
			}) {
				toDeleteGroups = append(toDeleteGroups, g.ID)
			}
		}

		if len(toDeleteGroups) > 0 {
			if res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE project_id = ? AND id IN(?)", (declarations.Group{}).TableName()), c.model.ProjectID, toDeleteGroups); res.Error != nil {
				return res.Error
			}
		}

		if len(toCreateGroups) > 0 {
			if res := tx.Create(&toCreateGroups); res.Error != nil {
				return res.Error
			}
		}

		return nil
	}); transactionErr != nil {
		return []declarations.Group{}, appErrors.NewApplicationError(transactionErr)
	}

	var groups []declarations.Group
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT name, id FROM %s WHERE project_id = ?", (declarations.Group{}).TableName()), c.model.ProjectID).Scan(&groups)
	if res.Error != nil {
		return []declarations.Group{}, appErrors.NewApplicationError(res.Error)
	}

	return groups, nil
}

func (c Main) Handle() ([]View, error) {
	if err := c.Validate(); err != nil {
		return []View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return []View{}, err
	}

	if err := c.Authorize(); err != nil {
		return []View{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return []View{}, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, []View, []declarations.Group] {
	logBuilder.Add("addToList", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
