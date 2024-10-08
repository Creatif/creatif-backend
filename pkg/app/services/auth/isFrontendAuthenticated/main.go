package isFrontendAuthenticated

import (
	"creatif/pkg/app/auth"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
)

type Main struct {
	auth auth.Authentication
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

func New(auth auth.Authentication) pkg.Job[interface{}, interface{}, interface{}] {
	return Main{auth: auth}
}
