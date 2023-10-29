package auth

import "time"

type apiAuthentication struct {
}

func (a *apiAuthentication) Authenticate() error {
	return nil
}

func (a *apiAuthentication) User() AuthenticatedUser {
	return AuthenticatedUser{
		ID:        "",
		Name:      "",
		LastName:  "",
		Email:     "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}

func NewApiAuthentication() Authentication {
	return &apiAuthentication{}
}
