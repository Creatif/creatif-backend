package getStructures

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/storage"
	"fmt"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return publicApiError.NewError("getStructures", errs, publicApiError.ValidationError)
	}

	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return publicApiError.NewError("getStructures", map[string]string{
			"unauthorized": "You are unauthorized to use this route",
		}, 403)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (LogicModel, error) {
	var lists []declarations.List
	var maps []declarations.Map

	if res := storage.Gorm().Raw(fmt.Sprintf("SELECT id, short_id, name, created_at, updated_at from %s WHERE project_id = ?", (declarations.List{}).TableName()), c.model.ProjectID).Scan(&lists); res.Error != nil {
		return LogicModel{}, publicApiError.NewError("getStructures", map[string]string{
			"dbFailed": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res := storage.Gorm().Raw(fmt.Sprintf("SELECT id, short_id, name, created_at, updated_at from %s WHERE project_id = ?", (declarations.Map{}).TableName()), c.model.ProjectID).Scan(&maps); res.Error != nil {
		return LogicModel{}, publicApiError.NewError("getStructures", map[string]string{
			"dbFailed": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return LogicModel{
		Maps:  maps,
		Lists: lists,
	}, nil
}

func (c Main) Handle() ([]View, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	if err := c.Authorize(); err != nil {
		return nil, err
	}

	model, err := c.Logic()

	if err != nil {
		return nil, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication) pkg.Job[Model, []View, LogicModel] {
	return Main{model: model, auth: auth}
}
