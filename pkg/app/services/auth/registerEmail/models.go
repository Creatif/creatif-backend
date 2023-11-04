package registerEmail

import (
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Model struct {
	Name           string
	LastName       string
	Email          string
	Password       string
	PolicyAccepted bool
}

func NewModel(name, lastName, email, password string, policyAccepted bool) Model {
	return Model{
		Name:           name,
		LastName:       lastName,
		Email:          email,
		Password:       password,
		PolicyAccepted: policyAccepted,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":           a.Name,
		"lastName":       a.LastName,
		"email":          a.Email,
		"policyAccepted": a.PolicyAccepted,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("lastName", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("email", is.Email),
			validation.Key("policyAccepted", validation.Required, validation.By(func(value interface{}) error {
				policyAccepted := value.(bool)

				if !policyAccepted {
					return errors.New("You cannot register a new account if you don't accept our policy.")
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
