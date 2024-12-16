package paginateListItems

import (
	"creatif/pkg/app/services/publicApi/publicApiError"
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/tidwall/sjson"
	"gorm.io/datatypes"
	"time"
)

type MarshalingConnectionItem struct {
	StructureID      string `json:"structureId"`
	StructureShortID string `json:"structureShortId"`
	StructureName    string `json:"structureName"`
	ProjectID        string `json:"projectId"`

	Name      string         `json:"name"`
	ID        string         `json:"id"`
	ShortID   string         `json:"shortId"`
	Value     datatypes.JSON `json:"value"`
	Behaviour string         `json:"behaviour"`
	Locale    string         `json:"locale"`
	Index     float64        `json:"index"`
	Groups    pq.StringArray `json:"groups"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type itemsConnections struct {
	item  *Item
	conns []ConnectionItem
}

func replaceConnectionJson(conns []ConnectionItem, items []Item, options Options) ([]Item, error) {
	itemConns := make([]itemsConnections, 0)
	for _, item := range items {
		c := make([]ConnectionItem, 0)
		for _, conn := range conns {
			if item.ItemID == conn.ParentVariableID {
				c = append(c, conn)
			}
		}

		itemConns = append(itemConns, itemsConnections{
			item:  &item,
			conns: c,
		})
	}

	for _, itemConn := range itemConns {
		item := itemConn.item
		updatedValue := item.Value

		for _, conn := range itemConn.conns {
			if options.ValueOnly {
				v, err := sjson.SetRawBytes(updatedValue, conn.Path, conn.Value)
				if err != nil {
					return nil, publicApiError.NewError("getMapItemByName", map[string]string{
						"internalError": err.Error(),
					}, publicApiError.DatabaseError)
				}
				updatedValue = v
			} else {
				b, err := json.Marshal(MarshalingConnectionItem{
					StructureID:      conn.StructureID,
					StructureShortID: conn.StructureShortID,
					StructureName:    conn.StructureName,
					ProjectID:        conn.ProjectID,
					Name:             conn.ItemName,
					ID:               conn.ItemID,
					ShortID:          conn.ItemShortID,
					Value:            conn.Value,
					Behaviour:        conn.Behaviour,
					Locale:           conn.Locale,
					Index:            conn.Index,
					Groups:           conn.Groups,
					CreatedAt:        conn.CreatedAt,
					UpdatedAt:        conn.UpdatedAt,
				})
				if err != nil {
					return nil, publicApiError.NewError("getMapItemByName", map[string]string{
						"internalError": err.Error(),
					}, publicApiError.DatabaseError)
				}

				v, err := sjson.SetRawBytes(updatedValue, conn.Path, b)
				if err != nil {
					return nil, publicApiError.NewError("getMapItemByName", map[string]string{
						"internalError": err.Error(),
					}, publicApiError.DatabaseError)
				}
				updatedValue = v
			}
		}

		item.Value = updatedValue
	}

	return sdk.Map(itemConns, func(idx int, value itemsConnections) Item {
		return *value.item
	}), nil
}
