package removeStructure

import "C"
import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Results struct {
	Errors []error
}

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("truncateStructure", "Validating...")
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

				referenceSql := fmt.Sprintf("DELETE FROM %s WHERE project_id = ? AND (parent_structure_id = ? OR child_structure_id = ?)", (declarations.Reference{}).TableName())
				if res := tx.Exec(referenceSql, c.model.ProjectID, c.model.ID, c.model.ID); res.Error != nil {
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

		referenceSql := fmt.Sprintf("DELETE FROM %s WHERE project_id = ? AND (parent_structure_id = ? OR child_structure_id = ?)", (declarations.Reference{}).TableName())
		if res := tx.Exec(referenceSql, c.model.ProjectID, c.model.ID, c.model.ID); res.Error != nil {
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, interface{}, interface{}] {
	logBuilder.Add("truncateStructure", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
