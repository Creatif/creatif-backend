package getVersions

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/published"
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
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
	c.logBuilder.Add("getVersions", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return publicApiError.NewError("getVersions", errs, publicApiError.ValidationError)
	}

	c.logBuilder.Add("getVersions", "Validated")
	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return publicApiError.NewError("getVersions", map[string]string{
			"unauthorized": "You are unauthorized to use this route",
		}, 403)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() ([]published.Version, error) {
	var version []published.Version
	if res := storage.Gorm().Raw(fmt.Sprintf("SELECT id, name, project_id, created_at, updated_at, is_production_version FROM %s WHERE project_id = ? ORDER BY created_at DESC", (published.Version{}).TableName()), c.model.ProjectID).Scan(&version); res.Error != nil {
		return []published.Version{}, publicApiError.NewError("getVersions", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return version, nil
}

func (c Main) Handle() ([]View, error) {
	if err := c.Validate(); err != nil {
		return []View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return []View{}, err
	}

	if err := c.Authorize(); err != nil {
		return []View{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return []View{}, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, []View, []published.Version] {
	logBuilder.Add("getVersions", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
