package auth

import "time"

type AuthenticatedUser struct {
	ID string

	Name     string
	LastName string
	Email    string

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type Authentication interface {
	Authenticate() error
	User() AuthenticatedUser
}
