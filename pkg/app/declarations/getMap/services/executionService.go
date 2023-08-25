package services

import (
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type Node struct {
	ID string `json:"id" gorm:"primarykey"`

	Name      string         `json:"name" gorm:"index;uniqueIndex:unique_node"`
	Type      string         `json:"type"`
	Value     datatypes.JSON `json:"value"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups" gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON `json:"metadata"`

	CreatedAt time.Time `json:"createdAt" gorm:"<-:create"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func Execute(mapId string, strategy QueryStrategy) ([]Node, error) {
	var node []Node
	if res := storage.Gorm().Raw(strategy.GetQuery(), mapId).Scan(&node); res.Error != nil {
		return nil, appErrors.NewDatabaseError(res.Error).AddError("getMap.Services.Execute", nil)
	}

	return node, nil
}
