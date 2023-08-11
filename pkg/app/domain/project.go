package domain

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	ID string `gorm:"primarykey"`

	Name string

	UserID *string
	User   User

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *Project) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (Project) TableName() string {
	return PROJECT_TABLE
}

func (u Project) WithSchema() string {
	return fmt.Sprintf("%s.%s", "app", u.TableName())
}

func NewProject(name string) Project {
	return Project{
		Name: name,
	}
}
