package loginApi

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Model struct {
	Email     string
	Password  string
	Session   string
	ApiKey    string
	ProjectID string
}

func NewModel(email, password, session string) Model {
	return Model{
		Email:    email,
		Password: password,
		Session:  session,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"password": a.Password,
		"session":  a.Session,
		"email":    a.Email,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("email", is.Email),
			validation.Key("password", validation.Required, validation.RuneLength(8, 20)),
			validation.Key("session", validation.Required),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
