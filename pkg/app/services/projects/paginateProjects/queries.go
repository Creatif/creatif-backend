package paginateProjects

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type Variable struct {
	ID string `gorm:"primarykey"`

	Name      string         `gorm:"index;uniqueIndex:unique_variable"`
	Behaviour string         // readonly,modifiable
	Groups    pq.StringArray `gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON
	Value     datatypes.JSON
	ProjectID string

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
