package auth

import "time"

type noopAuthentication struct {
}

func (a *noopAuthentication) Authenticate() error {
	return nil
}

func (a *noopAuthentication) User() AuthenticatedUser {
	return AuthenticatedUser{
		ID:        "",
		Name:      "",
		LastName:  "",
		Email:     "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}

func (a *noopAuthentication) Refresh() (string, error) {
	return "", nil
}

func (a *noopAuthentication) Logout(cb func()) {
}

func (a *noopAuthentication) ShouldRefresh() bool {
	return false
}

func NewNoopAuthentication() Authentication {
	return &noopAuthentication{}
}
