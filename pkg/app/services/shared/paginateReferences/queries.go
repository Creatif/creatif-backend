package paginateReferences

import (
	"gorm.io/datatypes"
	"time"
)

type Variable struct {
	ID string `gorm:"primarykey"`

	Name      string `gorm:"index;uniqueIndex:unique_variable"`
	Behaviour string // readonly,modifiable
	Metadata  datatypes.JSON
	Value     datatypes.JSON
	ProjectID string

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
