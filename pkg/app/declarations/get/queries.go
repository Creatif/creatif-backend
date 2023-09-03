package get

import (
	"github.com/lib/pq"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
	"time"
)

type NodeWithValueQuery struct {
	ID ksuid.KSUID

	Name      string
	Behaviour string // readonly,modifiable
	Groups    pq.StringArray
	Metadata  datatypes.JSON
	Value     interface{}

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
