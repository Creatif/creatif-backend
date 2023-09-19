package assignments

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MapVariableValue struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	MapVariableID string                   `gorm:"type:text;check:length(id)=26"`
	MapVariable   declarations.MapVariable `gorm:"foreignKey:MapID"`

	Value datatypes.JSON `gorm:"type:jsonb"`
}

func NewMapVariableValue(mapVariableId string, value datatypes.JSON) MapVariableValue {
	return MapVariableValue{
		Value:         value,
		MapVariableID: mapVariableId,
	}
}

func (u *MapVariableValue) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (MapVariableValue) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.MAP_VARIABLE_VALUE)
}
