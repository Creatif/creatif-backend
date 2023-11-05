package app

import "github.com/microcosm-cc/bluemonday"

type RegisterEmail struct {
	Name           string `json:"name"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PolicyAccepted bool   `json:"policyAccepted"`
}

func SanitizeRegisterEmail(model RegisterEmail) RegisterEmail {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.Email = p.Sanitize(model.Email)
	model.Password = p.Sanitize(model.Password)
	model.LastName = p.Sanitize(model.LastName)

	return model
}
