package getValue

import (
	"creatif/pkg/lib/storage"
)

func queryValue(name string) (Node, error) {
	var node Node
	if res := storage.Gorm().Raw(`SELECT n.value FROM declarations.nodes AS n WHERE n.name = ?`, name).Scan(&node); res.Error != nil {
		return Node{}, res.Error
	}

	return node, nil
}
