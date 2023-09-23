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

	ProjectID string `gorm:"index;type:text;check:length(id)=26"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewMap(projectId, name string) Map {
	return Map{
		Name:      name,
		ProjectID: projectId,
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
	return fmt.Sprintf("%s.%s", "declarations", domain.VARIABLE_MAP)
}
