package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       string `gorm:"primarykey"`
	Name     string
	LastName string
	Email    string

	Projects []Project

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (User) TableName() string {
	return USERS_TABLE
}

func NewUser(name string, lastName string, email string, isSuperAdmin bool, isAdmin bool) User {
	return User{
		Name:     name,
		LastName: lastName,
		Email:    email,
	}
}
