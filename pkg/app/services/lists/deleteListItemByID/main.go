package deleteListItemByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
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
	c.logBuilder.Add("deleteListItemByID", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("deleteListItemByID", "Validated")
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
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("deleteListItemByID", err.Error())
		return nil, appErrors.NewApplicationError(err).AddError("deleteListItemByID.Logic", nil)
	}

	listId, listVal := shared.DetermineID("l", c.model.Name, c.model.ID, c.model.ShortID)
	listVarId, listVarVal := shared.DetermineID("lv", "", c.model.ItemID, c.model.ItemShortID)
	sql := fmt.Sprintf(
		`DELETE FROM %s AS lv USING %s AS l WHERE %s AND l.project_id = ? AND l.locale_id = ? AND lv.list_id = l.id AND %s`,
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName(),
		listId,
		listVarId,
	)

	res := storage.Gorm().Exec(sql, listVal, c.model.ProjectID, localeID, listVarVal)
	if res.Error != nil {
		c.logBuilder.Add("deleteListItemByID", res.Error.Error())
		return nil, appErrors.NewDatabaseError(res.Error).AddError("deleteListItemByID.Logic", nil)
	}

	if res.RowsAffected == 0 {
		return nil, appErrors.NewNotFoundError(errors.New("List or variable not found")).AddError("deleteListItemByID.Logic", nil)
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
	logBuilder.Add("deleteListItemByID", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
