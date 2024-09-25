package deleteListItemByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/events"
	"creatif/pkg/app/services/shared"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	err := shared.IsParent(c.model.ItemID)
	if err != nil {
		return appErrors.NewValidationError(map[string]string{
			"isParent": "This variable is a parent and cannot be deleted",
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

func (c Main) Logic() (*struct{}, error) {
	if transactionErr := storage.Transaction(func(tx *gorm.DB) error {
		paths, err := getImagePaths(c.model.ProjectID, c.model.ItemID)
		if err != nil {
			return err
		}

		deleteImagesSql := fmt.Sprintf(
			`DELETE FROM %s WHERE list_id = ? AND project_id = ?`,
			(declarations.File{}).TableName(),
		)

		res := tx.Exec(deleteImagesSql, c.model.ItemID, c.model.ProjectID)
		if res.Error != nil {
			return appErrors.NewDatabaseError(res.Error).AddError("deleteListItemByID.Logic", nil)
		}

		sql := fmt.Sprintf(
			`DELETE FROM %s AS lv USING %s AS l WHERE (l.id = ? OR l.short_id = ?) AND l.project_id = ? AND lv.list_id = l.id AND (lv.id = ? OR lv.short_id = ?)`,
			(declarations.ListVariable{}).TableName(),
			(declarations.List{}).TableName(),
		)

		res = tx.Exec(sql, c.model.Name, c.model.Name, c.model.ProjectID, c.model.ItemID, c.model.ItemID)
		if res.Error != nil {
			return appErrors.NewDatabaseError(res.Error).AddError("deleteListItemByID.Logic", nil)
		}

		if res.RowsAffected == 0 {
			return appErrors.NewNotFoundError(errors.New("List or variable not found")).AddError("deleteListItemByID.Logic", nil)
		}

		if res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE variable_id = ?", (declarations.VariableGroup{}).TableName()), c.model.ItemID); res.Error != nil {
			return res.Error
		}

		if err := shared.RemoveAsParent(c.model.ItemID, tx); err != nil {
			return err
		}
		if err := shared.RemoveAsChild(c.model.ItemID, tx); err != nil {
			return err
		}

		for _, path := range paths {
			if err := os.Remove(path); err != nil {
				events.DispatchEvent(events.NewFileNotRemoveEvent(path, "", c.model.ProjectID))
			}
		}

		return nil
	}); transactionErr != nil {
		return nil, transactionErr
	}

	return nil, nil
}

func (c Main) Handle() (*struct{}, error) {
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, *struct{}, *struct{}] {
	return Main{model: model, auth: auth}
}
