package auth

import "time"

type frontendAuthentication struct {
}

func (a *frontendAuthentication) Authenticate() error {
	return nil
}

func (a *frontendAuthentication) User() AuthenticatedUser {
	return AuthenticatedUser{
		ID:        "",
		Name:      "",
		LastName:  "",
		Email:     "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}

func NewFrontendAuthentication() Authentication {
	return &frontendAuthentication{}
}
