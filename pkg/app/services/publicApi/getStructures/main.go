package getStructures

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/publicApi/publicApiError"
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
	c.logBuilder.Add("getStructures", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("getStructures", "Validated")
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

func (c Main) Logic() (LogicModel, error) {
	var lists []declarations.List
	var maps []declarations.Map

	if res := storage.Gorm().Raw(fmt.Sprintf("SELECT * from %s WHERE project_id = ?", (declarations.List{}).TableName()), c.model.ProjectID).Scan(&lists); res.Error != nil {
		return LogicModel{}, publicApiError.NewError("getStructures", map[string]string{
			"dbFailed": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res := storage.Gorm().Raw(fmt.Sprintf("SELECT * from %s WHERE project_id = ?", (declarations.Map{}).TableName()), c.model.ProjectID).Scan(&maps); res.Error != nil {
		return LogicModel{}, publicApiError.NewError("getStructures", map[string]string{
			"dbFailed": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return LogicModel{
		Maps:  maps,
		Lists: lists,
	}, nil
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, LogicModel] {
	logBuilder.Add("getStructures", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
