package removeMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/events"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
	"os"
)

type Main struct {
	model Model
	auth  auth.Authentication
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
	if err := storage.Transaction(func(tx *gorm.DB) error {
		paths, err := getImagePaths(c.model.ProjectID, c.model.VariableName)
		if err != nil {
			return err
		}

		deleteListImages := fmt.Sprintf(
			`DELETE FROM %s WHERE (list_id = ? OR map_id = ?) AND project_id = ?`,
			(declarations.File{}).TableName(),
		)

		res := tx.Exec(deleteListImages, c.model.VariableName, c.model.VariableName, c.model.ProjectID)
		if res.Error != nil {
			return appErrors.NewDatabaseError(res.Error).AddError("deleteListItemByID.Logic", nil)
		}

		sql := fmt.Sprintf(
			`DELETE FROM %s AS mv USING %s AS m WHERE m.project_id = ? AND mv.map_id = m.id AND (mv.id = ? OR mv.short_id = ?) AND (m.id = ? OR m.short_id = ?)`,
			(declarations.MapVariable{}).TableName(),
			(declarations.Map{}).TableName(),
		)

		res = tx.Exec(sql, c.model.ProjectID, c.model.VariableName, c.model.VariableName, c.model.Name, c.model.Name)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return res.Error
		}

		if res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE variable_id = ?", (declarations.VariableGroup{}).TableName()), c.model.VariableName); res.Error != nil {
			return res.Error
		}

		if res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE parent_variable_id = ?", (declarations.Connection{}).TableName()), c.model.VariableName); res.Error != nil {
			return res.Error
		}

		for _, path := range paths {
			if err := os.Remove(path); err != nil {
				events.DispatchEvent(events.NewFileNotRemoveEvent(path, "", c.model.ProjectID))
			}
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, interface{}, interface{}] {
	return Main{model: model, auth: auth}
}
