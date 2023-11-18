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

func NewModel(email, password, apiKey, projectId, session string) Model {
	return Model{
		Email:     email,
		Password:  password,
		ApiKey:    apiKey,
		Session:   session,
		ProjectID: projectId,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"email":     a.Email,
		"password":  a.Password,
		"apiKey":    a.ApiKey,
		"projectId": a.ProjectID,
		"session":   a.Session,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("email", is.Email),
			validation.Key("apiKey", validation.Required),
			validation.Key("password", validation.Required, validation.RuneLength(8, 20)),
			validation.Key("projectId", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("session", validation.Required),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
