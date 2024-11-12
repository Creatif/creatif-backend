package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/lib/pq"
)

type PublishedGroups struct {
	VariableID string         `gorm:"index;type:text"`
	VersionID  string         `gorm:"index:type:text"`
	ProjectID  string         `gorm:"index;type:text"`
	Groups     pq.StringArray `gorm:"type:text[];not_null"`
}

func (PublishedGroups) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.PUBLISHED_GROUPS_TABLE)
}
