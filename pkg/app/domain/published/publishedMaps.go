package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type PublishedMap struct {
	ID      string `gorm:"primaryKey;type:text"`
	ShortID string `gorm:"type:text"`

	VersionID string  `gorm:"primaryKey;type:text"`
	Version   Version `gorm:"foreignKey:VersionID"`

	Name string `gorm:"type:text"`

	VariableName    string         `gorm:"type:text"`
	VariableID      string         `gorm:"primaryKey;type:text"`
	VariableShortID string         `gorm:"type:text"`
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
	value datatypes.JSON,
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
