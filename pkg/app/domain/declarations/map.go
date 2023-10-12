package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Map struct {
	ID      string `gorm:"primarykey;type:text CHECK(length(id)=26)"`
	ShortID string `gorm:"uniqueIndex:unique_map;type:text"`

	Name string `gorm:"uniqueIndex:unique_map_name"`

	ProjectID    string        `gorm:"uniqueIndex:unique_map_name;type:text;check:length(id)=26;not null;default: null"`
	MapVariables []MapVariable `gorm:"foreignKey:MapID;constraint:OnDelete:CASCADE;"`

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
	shortId, err := storage.ShortId.Generate()
	if err != nil {
		return err
	}
	u.ShortID = shortId

	return nil
}

func (Map) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.VARIABLE_MAP)
}
