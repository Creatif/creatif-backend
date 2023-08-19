package create

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type NodeWithValueQuery struct {
	ID string

	Name      string
	Type      string // text,image,file,boolean
	Behaviour string // readonly,modifiable
	Groups    pq.StringArray
	Metadata  datatypes.JSON
	Value     interface{}

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}