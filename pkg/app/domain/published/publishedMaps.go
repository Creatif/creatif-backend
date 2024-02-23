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

	VersionID string `gorm:"type:text"`

	Name string `gorm:"index;type:text"`

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

func NewPublishedMap(
	id,
	shortId,
	versionId,
	name,
	variableName,
	variableId,
	variableShortId,
	behaviour,
	locale string,
	value []byte,
	groups []string,
	index float64,
) PublishedMap {
	return PublishedMap{
		Name:            name,
		ID:              id,
		ShortID:         shortId,
		VersionID:       versionId,
		VariableName:    variableName,
		VariableID:      variableId,
		VariableShortID: variableShortId,
		Behaviour:       behaviour,
		LocaleID:        locale,
		Value:           value,
		Groups:          groups,
		Index:           index,
	}
}

func (PublishedMap) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.PUBLISHED_MAPS_TABLE)
}
