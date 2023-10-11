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

	Name    string
	Alpha3b string
	Alpha3t string
	Alpha2  string

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewLanguage(name, alpha3b, alpha3t, alpha2 string) Language {
	return Language{
		Name:    name,
		Alpha3b: alpha3b,
		Alpha3t: alpha3t,
		Alpha2:  alpha2,
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
