package services

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
)

func GetMap(id string) (declarations.Map, error) {
	var m declarations.Map
	if sdk.IsValidUuid(id) {
		if err := storage.Get(m.TableName(), id, &m, "ID", "Name", "CreatedAt", "UpdatedAt"); err != nil {
			return declarations.Map{}, appErrors.NewDatabaseError(err).AddError("Map.Get.Logic", nil)
		}
	} else {
		if err := storage.GetBy(m.TableName(), "name", id, &m, "ID", "Name", "CreatedAt", "UpdatedAt"); err != nil {
			return declarations.Map{}, appErrors.NewDatabaseError(err).AddError("Map.Get.Logic", nil)
		}
	}

	return m, nil
}
