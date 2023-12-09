package deleteRangeByID

import "C"
import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
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
	c.logBuilder.Add("deleteRangeByID", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("deleteRangeByID", "Validated")
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

func (c Main) Logic() (*struct{}, error) {
	listId, listVal := shared.DetermineID("l", c.model.Name, c.model.ID, c.model.ShortID)
	sql := fmt.Sprintf(
		`DELETE FROM %s AS lv USING %s AS l WHERE %s AND l.project_id = ? AND lv.list_id = l.id AND lv.id IN(?)`,
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName(),
		listId,
	)

	res := storage.Gorm().Exec(sql, listVal, c.model.ProjectID, c.model.Items)
	if res.Error != nil {
		c.logBuilder.Add("deleteRangeByID", res.Error.Error())
		return nil, appErrors.NewDatabaseError(res.Error).AddError("deleteRangeByID.Logic", nil)
	}

	if res.RowsAffected == 0 {
		return nil, appErrors.NewNotFoundError(res.Error).AddError("deleteRangeByID.Logic", nil)
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, *struct{}, *struct{}] {
	logBuilder.Add("deleteRangeByID", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
