package createAdmin

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Model struct {
	Name     string
	LastName string
	Email    string
	Password string
}

func NewModel(name, lastName, email, password string) Model {
	return Model{
		Name:     name,
		LastName: lastName,
		Email:    email,
		Password: password,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":     a.Name,
		"lastName": a.LastName,
		"email":    a.Email,
		"password": a.Password,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("lastName", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("email", is.EmailFormat, is.Email),
			validation.Key("password", validation.Required, validation.RuneLength(8, 20)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
