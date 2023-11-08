package auth

import "time"

type AuthenticatedUser struct {
	ID string `json:"id"`

	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`

	Refresh time.Time `json:"refresh"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AuthenticationSession struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func NewAuthenticatedUser(id, name, lastName, email string, createdAt, updatedAt, refresh time.Time) AuthenticatedUser {
	return AuthenticatedUser{
		ID:       id,
		Name:     name,
		LastName: lastName,
		Email:    email,

		Refresh: refresh,

		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewAuthenticationSession(id, token string) AuthenticationSession {
	return AuthenticationSession{
		ID:    id,
		Token: token,
	}
}

type Authentication interface {
	Authenticate() error
	User() AuthenticatedUser
	ShouldRefresh() bool
	Refresh() (string, error)
}

type Loginer interface {
	Login() (string, error)
}
