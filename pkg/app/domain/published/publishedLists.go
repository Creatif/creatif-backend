package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type PublishedList struct {
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

func NewPublishedList(
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
) PublishedList {
	return PublishedList{
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

func (PublishedList) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.PUBLISHED_LISTS_TABLE)
}
