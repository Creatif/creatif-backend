package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Map struct {
	ID      string `gorm:"primarykey;type:text;default:gen_ulid()"`
	ShortID string `gorm:"uniqueIndex:unique_map;type:text"`

	Name string `gorm:"uniqueIndex:unique_map_name;not null"`

	ProjectID    string        `gorm:"uniqueIndex:unique_map_name;type:text"`
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
