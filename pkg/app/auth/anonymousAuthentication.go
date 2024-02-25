package auth

type anonymousAuthentication struct {
}

func (a *anonymousAuthentication) Authenticate() error {
	return nil
}

func (a *anonymousAuthentication) User() AuthenticatedUser {
	return AuthenticatedUser{}
}

func (a *anonymousAuthentication) ShouldRefresh() bool {
	return false
}

func (a *anonymousAuthentication) Logout(cb func()) {
	cb()
}

func (a *anonymousAuthentication) Refresh() (string, error) {
	return "", nil
}

func NewAnonymousAuthentication() Authentication {
	return &anonymousAuthentication{}
}
