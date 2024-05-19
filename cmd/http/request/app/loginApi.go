package app

import "github.com/microcosm-cc/bluemonday"

type LoginApi struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SanitizeLoginApi(model LoginApi) LoginApi {
	p := bluemonday.StrictPolicy()
	model.Email = p.Sanitize(model.Email)
	model.Password = p.Sanitize(model.Password)

	return model
}
