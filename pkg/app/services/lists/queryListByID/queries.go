package queryListByID

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type QueryVariable struct {
	ID      string  `gorm:"primarykey;type:text"`
	ShortID string  `gorm:"uniqueIndex:unique_map_variable;type:text;not null"`
	Index   float64 `gorm:"type:float"`

	Name      string         `gorm:"uniqueIndex:unique_map_variable;not null"`
	Behaviour string         `gorm:"not null"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	Value     datatypes.JSON `gorm:"type:jsonb"`
	Groups    pq.StringArray `gorm:"type:text[];column:groups"`

	MapID    string `gorm:"uniqueIndex:unique_map_variable;type:text"`
	LocaleID string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
}
