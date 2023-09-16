package pagination

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type NodeWithValue struct {
	ID string `gorm:"primarykey;type:char(27)"`

	Name      string         `gorm:"index;uniqueIndex:unique_node"`
	Behaviour string         // readonly,modifiable
	Groups    pq.StringArray `gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON
	Value     datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NodeWithoutValue struct {
	ID string `gorm:"primarykey"`

	Name      string         `gorm:"index;uniqueIndex:unique_node"`
	Behaviour string         // readonly,modifiable
	Groups    pq.StringArray `gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
