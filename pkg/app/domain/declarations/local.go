package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Locale struct {
	ID string `gorm:"primarykey;type:text;default:gen_ulid()"`

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
	return nil
}

func (Locale) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.LOCALE_TABLE)
}
