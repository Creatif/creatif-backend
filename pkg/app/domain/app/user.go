package app

import (
	"creatif/pkg/app/domain"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID string `gorm:"primarykey;type:text;default:gen_ulid()"`

	Name     string
	LastName string
	Email    string `gorm:"index"`
	Password string

	Confirmed      bool
	PolicyAccepted bool

	Provider string

	Projects []Project `gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewUser(name, lastName, email, password, provider string, confirmed, policyAccepted bool) User {
	return User{
		Name:           name,
		LastName:       lastName,
		Email:          email,
		Password:       password,
		Provider:       provider,
		Confirmed:      confirmed,
		PolicyAccepted: policyAccepted,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (User) TableName() string {
	return fmt.Sprintf("%s.%s", "app", domain.USERS_TABLE)
}
