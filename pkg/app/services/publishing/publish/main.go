package publish

import "C"
import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/published"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
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
	c.logBuilder.Add("publish", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("publish", "Validated")
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

func (c Main) Logic() (published.Version, error) {
	version := published.NewVersion(c.model.ProjectID)
	if transactionError := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&version); res.Error != nil {
			return res.Error
		}

		fmt.Println(version.ID)

		listSql := fmt.Sprintf(`
	MERGE INTO %s AS p
	USING (%s) AS t
	ON p.variable_id != t.variableId
	WHEN NOT MATCHED THEN
        INSERT (
			id, 
			short_id, 
			version_id, 
			name, 
			variable_name, 
			variable_id, 
			variable_short_id, 
			index, 
			behaviour, 
			value, 
			locale_id, 
			groups
		) VALUES (
			t.ID, 
			t.shortId, 
			'%s', 
			t.name, 
			t.variableName, 
			t.variableId, 
			t.variableShortId, 
			t.index, 
			t.behaviour, 
			t.value, 
			t.locale, 
			t.groups
		)
`,
			(published.PublishedList{}).TableName(),
			getSelectListSql(c.model.ProjectID),
			version.ID,
		)

		mapSql := fmt.Sprintf(`
	MERGE INTO %s
	USING (%s) AS t
	ON p.variable_id != t.variableId
	WHEN NOT MATCHED THEN
        INSERT (
			id, 
			short_id, 
			version_id, 
			name, 
			variable_name, 
			variable_id, 
			variable_short_id, 
			index, 
			behaviour, 
			value, 
			locale_id, 
			groups
		) VALUES (
			t.ID, 
			t.shortId, 
			'%s', 
			t.name, 
			t.variableName, 
			t.variableId, 
			t.variableShortId, 
			t.index, 
			t.behaviour, 
			t.value, 
			t.locale, 
			t.groups
		)
`,
			(published.PublishedList{}).TableName(),
			getSelectMapSql(c.model.ProjectID),
			version.ID,
		)

		if res := tx.Exec(listSql, c.model.ProjectID); res.Error != nil {
			return res.Error
		}

		if res := tx.Exec(mapSql, c.model.ProjectID); res.Error != nil {
			return res.Error
		}

		return nil
	}); transactionError != nil {
		return published.Version{}, appErrors.NewApplicationError(transactionError)
	}

	return version, nil
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, published.Version] {
	logBuilder.Add("deleteRangeByID", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
