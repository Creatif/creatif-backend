package auth

import "time"

type AuthenticatedUser struct {
	ID string `json:"id"`

	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`

	Refresh   time.Time `json:"refresh"`
	ProjectID string    `json:"projectID"`
	ApiKey    string    `json:"apiKey"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AuthenticatedFrontendSession struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Type  string `json:"type"`
}

type AuthenticatedApiSession struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Type  string `json:"type"`
}

func NewAuthenticatedUser(id, name, lastName, email string, createdAt, updatedAt, refresh time.Time, projectID, apiKey string) AuthenticatedUser {
	return AuthenticatedUser{
		ID:        id,
		Name:      name,
		LastName:  lastName,
		Email:     email,
		ApiKey:    apiKey,
		ProjectID: projectID,

		Refresh: refresh,

		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewAuthenticatedFrontendSession(id, token string) AuthenticatedFrontendSession {
	return AuthenticatedFrontendSession{
		ID:    id,
		Token: token,
		Type:  "frontend",
	}
}

func NewAuthenticatedApiSession(id, token string) AuthenticatedApiSession {
	return AuthenticatedApiSession{
		ID:    id,
		Token: token,
		Type:  "api",
	}
}

type Authentication interface {
	Authenticate() error
	User() AuthenticatedUser
	ShouldRefresh() bool
	Refresh() (string, error)
	Logout(func())
}

type Loginer interface {
	Login() (string, error)
}
