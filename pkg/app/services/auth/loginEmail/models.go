package loginEmail

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Model struct {
	Email    string
	Password string
}

func NewModel(email, password string) Model {
	return Model{
		Email:    email,
		Password: password,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"email": a.Email,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("email", is.Email),
			validation.Key("email", validation.Required, validation.RuneLength(8, 20)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
