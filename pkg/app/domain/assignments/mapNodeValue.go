package assignments

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MapNodeValue struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	MapNodeID string               `gorm:"type:text;check:length(id)=26"`
	MapNode   declarations.MapNode `gorm:"foreignKey:MapID"`

	Value datatypes.JSON
}

func NewMapNodeValue(mapNodeId string, value datatypes.JSON) MapNodeValue {
	return MapNodeValue{
		Value:     value,
		MapNodeID: mapNodeId,
	}
}

func (u *MapNodeValue) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (MapNodeValue) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_MAP_VALUE_NODE)
}
