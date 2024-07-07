package removeVersion

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/published"
	"creatif/pkg/app/services/events"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"fmt"
	"os"
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
	c.logBuilder.Add("removeVersion", "Validating...")
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

func (c Main) Logic() (*struct{}, error) {
	var version published.Version
	if res := storage.Gorm().Raw(fmt.Sprintf("SELECT id, name FROM %s WHERE project_id = ? AND id = ?", (published.Version{}).TableName()), c.model.ProjectID, c.model.ID).Scan(&version); res.Error != nil {
		return nil, appErrors.NewApplicationError(res.Error)
	}

	if res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE project_id = ? AND id = ?", (published.Version{}).TableName()), c.model.ProjectID, c.model.ID); res.Error != nil {
		return nil, appErrors.NewApplicationError(res.Error)
	}

	path := fmt.Sprintf("%s/%s/%s", constants.PublicDirectory, c.model.ProjectID, version.Name)
	if err := os.RemoveAll(path); err != nil {
		events.DispatchEvent(events.NewPublicDirectoryNotRemoved(path, "", c.model.ProjectID))
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
	logBuilder.Add("removeVersion", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
