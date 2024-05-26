package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"time"
)

type Locale struct {
	ID string `gorm:"primarykey;type:text"`

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
	u.ID = ksuid.New().String()
	return nil
}

func (Locale) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.LOCALE_TABLE)
}
