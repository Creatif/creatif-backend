package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Language struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	Name  string
	Alpha string

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewLanguage(name, alpha string) Language {
	return Language{
		Name:  name,
		Alpha: alpha,
	}
}

func (u *Language) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (Language) TableName() string {
	return fmt.Sprintf("%s.%s", "app", domain.LANGUAGE_TABLE)
}
