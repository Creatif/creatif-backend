package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type PublishedMap struct {
	ID      string `gorm:"index;type:text"`
	ShortID string `gorm:"index;type:text"`
	Version string `gorm:"uniqueIndex;type:text"`

	Name      string `gorm:"index;type:text"`
	ProjectID string `gorm:"index;type:text"`

	VariableName    string         `gorm:"index;type:text"`
	VariableID      string         `gorm:"index;type:text"`
	VariableShortID string         `gorm:"index;type:text"`
	Index           float64        `gorm:"type:float"`
	Behaviour       string         `gorm:"not null"`
	Value           datatypes.JSON `gorm:"type:jsonb"`
	LocaleID        string         `gorm:"type:text"`
	Groups          pq.StringArray `gorm:"type:text[];not_null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPublishedMap(projectId, name string) PublishedMap {
	return PublishedMap{
		Name:      name,
		ProjectID: projectId,
	}
}

func (PublishedMap) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.PUBLISHED_MAPS_TABLE)
}
