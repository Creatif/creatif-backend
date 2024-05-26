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

	State string `gorm:"default: 'draft'"`

	Maps  []declarations.Map  `gorm:"foreignKey:ProjectID;references:ID;constraint:OnDelete:CASCADE;"`
	Lists []declarations.List `gorm:"foreignKey:ProjectID;references:ID;constraint:OnDelete:CASCADE;"`

	UserID string `gorm:"type:text CHECK(length(id)=26);not null"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewProject(name, userID string) Project {
	return Project{
		Name:   name,
		UserID: userID,
	}
}

func (u *Project) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (Project) TableName() string {
	return fmt.Sprintf("%s.%s", "app", domain.PROJECT_TABLE)
}
