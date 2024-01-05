package addToMap

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
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
	c.logBuilder.Add("addToMap", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}
	c.logBuilder.Add("addToMap", "Validated.")

	entryLocaleId, _ := locales.GetIDWithAlpha(c.model.Entry.Locale)

	sql := fmt.Sprintf(`
SELECT mv.id FROM %s AS mv 
INNER JOIN %s AS m ON 
(m.id = ? OR m.name = ? OR m.short_id = ?) AND m.project_id = ? AND 
mv.map_id = m.id AND mv.name = ? AND mv.locale_id = ?
`, (declarations.MapVariable{}).TableName(), (declarations.Map{}).TableName())

	var entry declarations.MapVariable
	res := storage.Gorm().Raw(sql, c.model.Name, c.model.Name, c.model.Name, c.model.ProjectID, c.model.Entry.Name, entryLocaleId).Scan(&entry)
	if res.Error != nil {
		return appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Variable with name '%s' and locale '%s' for map with ID '%s' already exists.", c.model.Entry.Name, c.model.Entry.Locale, c.model.Name),
		})
	}

	if res.RowsAffected != 0 {
		return appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Variable with name '%s' and locale '%s' for map with ID '%s' already exists.", c.model.Entry.Name, c.model.Entry.Locale, c.model.Name),
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

func (c Main) Logic() (interface{}, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Entry.Locale)
	if err != nil {
		c.logBuilder.Add("addToMap", err.Error())
		return nil, appErrors.NewApplicationError(err).AddError("addToMap.Logic", nil)
	}

	var m declarations.Map
	if res := storage.Gorm().Where(fmt.Sprintf("project_id = ? AND (id = ? OR name = ? OR short_id = ?)"), c.model.ProjectID, c.model.Name, c.model.Name, c.model.Name).Select("ID").First(&m); res.Error != nil {
		c.logBuilder.Add("addToMap", res.Error.Error())
		return nil, appErrors.NewNotFoundError(res.Error).AddError("addToMap.Logic", nil)
	}

	if c.model.Entry.Groups == nil {
		c.model.Entry.Groups = []string{}
	}

	mapNode := declarations.NewMapVariable(m.ID, localeID, c.model.Entry.Name, c.model.Entry.Behaviour, c.model.Entry.Metadata, c.model.Entry.Groups, c.model.Entry.Value)
	if res := storage.Gorm().Create(&mapNode); res.Error != nil {
		c.logBuilder.Add("addToMap", res.Error.Error())

		return nil, appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Map with name '%s' already exists.", c.model.Entry.Name),
		})
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
	logBuilder.Add("getMap", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
