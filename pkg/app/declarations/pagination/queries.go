package pagination

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type NodeWithValue struct {
	ID uuid.UUID `gorm:"primarykey"`

	Name      string         `gorm:"index;uniqueIndex:unique_node"`
	Behaviour string         // readonly,modifiable
	Groups    pq.StringArray `gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON
	Value     datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NodeWithoutValue struct {
	ID uuid.UUID `gorm:"primarykey"`

	Name      string         `gorm:"index;uniqueIndex:unique_node"`
	Behaviour string         // readonly,modifiable
	Groups    pq.StringArray `gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func getQueryWithoutValue() string {
	return `SELECT n.id, n.name, n.behaviour, n.metadata, n.groups, n.created_at, n.updated_at FROM declarations.nodes AS n`
}
