package app

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	Name string `gorm:"index"`

	Variables []declarations.Variable `gorm:"foreignKey:ProjectID;references:ID"`
	Maps      []declarations.Map      `gorm:"foreignKey:ProjectID;references:ID"`
	Lists     []declarations.List     `gorm:"foreignKey:ProjectID;references:ID"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewProject(name string) Project {
	return Project{
		Name: name,
	}
}

func (u *Project) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (Project) TableName() string {
	return fmt.Sprintf("%s.%s", "app", domain.PROJECT_TABLE)
}
