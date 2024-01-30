package auth

type frontendTestingAuthentication struct {
	authenticatedUser AuthenticatedUser
}

func (a *frontendTestingAuthentication) Authenticate() error {
	return nil
}

func (a *frontendTestingAuthentication) User() AuthenticatedUser {
	return a.authenticatedUser
}

func (a *frontendTestingAuthentication) Refresh() (string, error) {
	return "", nil
}

func (a *frontendTestingAuthentication) Logout(cb func()) {
}

func (a *frontendTestingAuthentication) ShouldRefresh() bool {
	return false
}

func NewFrontendTestingAuthentication(user AuthenticatedUser) Authentication {
	return &frontendTestingAuthentication{authenticatedUser: user}
}
