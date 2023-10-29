package app

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/app/domain/declarations"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	ID string `gorm:"primarykey;type:text;default:gen_ulid()"`

	Name string `gorm:"index"`

	Variables []declarations.Variable `gorm:"foreignKey:ProjectID;references:ID"`
	Maps      []declarations.Map      `gorm:"foreignKey:ProjectID;references:ID"`
	Lists     []declarations.List     `gorm:"foreignKey:ProjectID;references:ID"`

	UserID string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewProject(name string) Project {
	return Project{
		Name: name,
	}
}

func (u *Project) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (Project) TableName() string {
	return fmt.Sprintf("%s.%s", "app", domain.PROJECT_TABLE)
}
