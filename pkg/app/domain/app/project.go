package app

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/app/domain/declarations"
	"fmt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	ID string `gorm:"primarykey;type:text;default:gen_ulid()"`

	Name   string `gorm:"index"`
	APIKey string `gorm:"uniqueIndex"`

	Variables []declarations.Variable `gorm:"foreignKey:ProjectID;references:ID"`
	Maps      []declarations.Map      `gorm:"foreignKey:ProjectID;references:ID"`
	Lists     []declarations.List     `gorm:"foreignKey:ProjectID;references:ID"`

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
	id := ksuid.New().String()
	cost := 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(id), cost)
	if err != nil {
		return err
	}

	u.APIKey = string(bytes)
	return nil
}

func (Project) TableName() string {
	return fmt.Sprintf("%s.%s", "app", domain.PROJECT_TABLE)
}
