package paginateVariables

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type VariableWithValue struct {
	ID string `gorm:"primarykey;type:char(27)"`

	Name      string         `gorm:"index;uniqueIndex:unique_variable"`
	Behaviour string         // readonly,modifiable
	Groups    pq.StringArray `gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON
	Value     datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type VariableWithoutValue struct {
	ID string `gorm:"primarykey"`

	Name      string         `gorm:"index;uniqueIndex:unique_variable"`
	Behaviour string         // readonly,modifiable
	Groups    pq.StringArray `gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
