package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type MapVariable struct {
	ID      string `gorm:"primarykey;type:text CHECK(length(id)=26)"`
	ShortID string `gorm:"uniqueIndex:unique_variable;type:text"`

	Name      string `gorm:"uniqueIndex:unique_map_variable"`
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	Value     datatypes.JSON `gorm:"type:jsonb"`

	MapID string `gorm:"uniqueIndex:unique_map_variable;type:text;check:length(id)=26"`
	Map   Map    `gorm:"foreignKey:MapID"`

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
}

func NewMapVariable(mapId, name, behaviour string, metadata datatypes.JSON, groups pq.StringArray, value datatypes.JSON) MapVariable {
	return MapVariable{
		MapID:     mapId,
		Name:      name,
		Behaviour: behaviour,
		Metadata:  metadata,
		Groups:    groups,
		Value:     value,
	}
}

func (u *MapVariable) BeforeCreate(tx *gorm.DB) (err error) {
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

func (MapVariable) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.MAP_VARIABLES)
}
