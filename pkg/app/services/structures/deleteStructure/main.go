package removeStructure

import "C"
import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Results struct {
	Errors []error
}

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
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

func (c Main) Logic() (interface{}, error) {
	if c.model.Type == "map" {
		if c.model.Type == "map" {
			if err := storage.Gorm().Transaction(func(tx *gorm.DB) error {
				mapSql := fmt.Sprintf("DELETE FROM %s WHERE id = ? AND project_id = ?", (declarations.Map{}).TableName())
				if res := tx.Exec(mapSql, c.model.ID, c.model.ProjectID); res.Error != nil {
					return errors.New(res.Error.Error())
				}

				connectionsSql := fmt.Sprintf(
					"DELETE FROM %s AS c USING %s AS v, %s AS vg WHERE c.project_id = ? AND (vg.id = child_variable_id OR vg.id = parent_variable_id) AND vg.map_id = v.id AND (c.parent_variable_id = ? OR c.child_variable_id = ?)",
					(declarations.Connection{}).TableName(),
					(declarations.Map{}).TableName(),
					(declarations.MapVariable{}).TableName(),
				)

				if res := tx.Exec(connectionsSql, c.model.ProjectID, c.model.ID, c.model.ID); res.Error != nil {
					return errors.New(res.Error.Error())
				}

				return nil
			}); err != nil {
				return nil, appErrors.NewDatabaseError(err)
			}

		}
	}

	if err := storage.Gorm().Transaction(func(tx *gorm.DB) error {
		mapSql := fmt.Sprintf("DELETE FROM %s WHERE id = ? AND project_id = ?", (declarations.List{}).TableName())
		if res := tx.Exec(mapSql, c.model.ID, c.model.ProjectID); res.Error != nil {
			return errors.New(res.Error.Error())
		}

		connectionsSql := fmt.Sprintf(
			"DELETE FROM %s AS c USING %s AS v, %s AS vg WHERE c.project_id = ? AND (vg.id = child_variable_id OR vg.id = parent_variable_id) AND vg.list_id = v.id AND (c.parent_variable_id = ? OR c.child_variable_id = ?)",
			(declarations.Connection{}).TableName(),
			(declarations.List{}).TableName(),
			(declarations.ListVariable{}).TableName(),
		)

		if res := tx.Exec(connectionsSql, c.model.ProjectID, c.model.ID, c.model.ID); res.Error != nil {
			return errors.New(res.Error.Error())
		}

		return nil
	}); err != nil {
		return nil, appErrors.NewDatabaseError(err)
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
