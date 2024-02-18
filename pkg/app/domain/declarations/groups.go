package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"time"
)

type Group struct {
	ID string `gorm:"primarykey;type:text;default:gen_ulid()"`

	Name      string `gorm:"uniqueIndex:unique_group;type:text"`
	ProjectID string `gorm:"uniqueIndex:unique_group;type:text"`

	VariableGroups []VariableGroup `gorm:"foreignKey:GroupID;references:ID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time `gorm:"<-:update"`
}

func NewGroup(projectId, name string) Group {
	return Group{
		Name:      name,
		ProjectID: projectId,
	}
}

func (Group) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.GROUPS_TABLE)
}
