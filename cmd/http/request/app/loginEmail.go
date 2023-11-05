package app

import "github.com/microcosm-cc/bluemonday"

type LoginEmail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SanitizeLoginEmail(model LoginEmail) LoginEmail {
	p := bluemonday.StrictPolicy()
	model.Email = p.Sanitize(model.Email)
	model.Password = p.Sanitize(model.Password)

	return model
}
