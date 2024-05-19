package createAndDiff

import "C"
import (
	"creatif/pkg/app/auth"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
)

type Results struct {
	Errors []error
}

type Main struct {
	model      interface{}
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
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

func New(model interface{}, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[interface{}, interface{}, interface{}] {
	logBuilder.Add("publish", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
