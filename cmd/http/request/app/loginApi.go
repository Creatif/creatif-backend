package app

import "github.com/microcosm-cc/bluemonday"

type LoginApi struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ProjectID string `json:"projectID"`
	ApiKey    string `json:"apiKey"`
	Session   string `json:"session"`
}

func SanitizeLoginApi(model LoginApi) LoginApi {
	p := bluemonday.StrictPolicy()
	model.Email = p.Sanitize(model.Email)
	model.Password = p.Sanitize(model.Password)
	model.Session = p.Sanitize(model.Session)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ApiKey = p.Sanitize(model.ApiKey)

	return model
}
