package get

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type NodeWithValueQuery struct {
	ID string

	Name      string
	Behaviour string // readonly,modifiable
	Groups    pq.StringArray
	Metadata  datatypes.JSON
	Value     interface{}

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
