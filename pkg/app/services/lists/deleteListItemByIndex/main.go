package deleteListItemByIndex

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
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

func (c Main) Logic() (*struct{}, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		return nil, appErrors.NewApplicationError(err).AddError("deleteListItemByIndex.Logic", nil)
	}

	offset := c.model.ItemIndex
	var existing declarations.ListVariable
	if res := storage.Gorm().
		Raw(fmt.Sprintf(`
			SELECT lv.id
			FROM %s AS lv INNER JOIN %s AS l
			ON l.project_id = ? AND l.name = ? AND lv.list_id = l.id AND l.locale_id = ?
			OFFSET ? LIMIT 1`, (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()), c.model.ProjectID, c.model.Name, localeID, offset).
		Scan(&existing); res.Error != nil || res.RowsAffected == 0 {
		if res.RowsAffected == 0 {
			return nil, appErrors.NewNotFoundError(res.Error).AddError("deleteListItemByIndex.Logic", nil)
		}

		return nil, appErrors.NewDatabaseError(res.Error).AddError("deleteListItemByIndex.Logic", nil)
	}

	var variable declarations.ListVariable
	res := storage.Gorm().Where("id = ? AND locale_id = ?", existing.ID, localeID).Delete(&variable)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError(res.Error).AddError("deleteListItemByIndex.Logic", nil)
		}

		return nil, appErrors.NewDatabaseError(res.Error).AddError("deleteListItemByIndex.Logic", nil)
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

func New(model Model) pkg.Job[Model, *struct{}, *struct{}] {
	return Main{model: model}
}
