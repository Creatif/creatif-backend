package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type MapVariable struct {
	ID      string  `gorm:"primarykey;type:text;default:gen_ulid()"`
	ShortID string  `gorm:"uniqueIndex:unique_map_variable;type:text;not null"`
	Index   float64 `gorm:"type:float"`

	Name      string         `gorm:"uniqueIndex:unique_map_variable;not null"`
	Behaviour string         `gorm:"not null"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	Value     datatypes.JSON `gorm:"type:jsonb"`

	MapID    string `gorm:"uniqueIndex:unique_map_variable;type:text"`
	LocaleID string `gorm:"type:text"`
	Map      Map    `gorm:"foreignKey:MapID"`

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
}

func NewMapVariable(mapId, localeID, name, behaviour string, metadata datatypes.JSON, value datatypes.JSON) MapVariable {
	return MapVariable{
		MapID:     mapId,
		LocaleID:  localeID,
		Name:      name,
		Behaviour: behaviour,
		Metadata:  metadata,
		Value:     value,
	}
}

func (u *MapVariable) BeforeCreate(tx *gorm.DB) (err error) {
	shortId, err := storage.ShortId.Generate()
	if err != nil {
		return err
	}
	u.ShortID = shortId

	return nil
}

func (MapVariable) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.MAP_VARIABLES)
}
