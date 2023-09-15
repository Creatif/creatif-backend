package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Map struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	Name string `gorm:"uniqueIndex"`

	CreatedAt time.Time `gorm:"<-:createNode"`
	UpdatedAt time.Time
}

func NewMap(name string) Map {
	return Map{
		Name: name,
	}
}

func (u *Map) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (Map) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.NODE_MAP_TABLE)
}
