package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Locale struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	Name  string
	Alpha string `gorm:"unique"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewLocale(name, alpha string) Locale {
	return Locale{
		Name:  name,
		Alpha: alpha,
	}
}

func (u *Locale) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (Locale) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.LOCALE_TABLE)
}
