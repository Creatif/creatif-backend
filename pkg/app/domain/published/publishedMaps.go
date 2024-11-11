package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type PublishedMap struct {
	StructureID string `gorm:"primaryKey;type:text"`
	ShortID     string `gorm:"primaryKey;type:text"`

	VersionID string  `gorm:"primaryKey;type:text"`
	Version   Version `gorm:"foreignKey:VersionID"`

	Name string `gorm:"primaryKey;type:text"`

	VariableName    string         `gorm:"primaryKey;type:text"`
	VariableID      string         `gorm:"primaryKey;type:text"`
	VariableShortID string         `gorm:"primaryKey;type:text"`
	Index           float64        `gorm:"type:float"`
	Behaviour       string         `gorm:"not null"`
	Value           datatypes.JSON `gorm:"type:jsonb"`
	LocaleID        string         `gorm:"type:text"`
	Groups          pq.StringArray `gorm:"type:text[];not_null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (PublishedMap) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.PUBLISHED_MAPS_TABLE)
}
