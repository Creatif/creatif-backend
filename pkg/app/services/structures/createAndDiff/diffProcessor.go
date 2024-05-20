package createAndDiff

import "creatif/pkg/app/domain/declarations"

func processListsAndMaps(
	projectId string,
	configStructures []Structure,
	lists []declarations.List,
	maps []declarations.Map,
) ([]declarations.List, []declarations.Map, []declarations.List, []declarations.Map) {
	listsToCreate := make([]declarations.List, 0)
	mapsToCreate := make([]declarations.Map, 0)

	listExistsInDb := make([]declarations.List, 0)
	mapExistsInDb := make([]declarations.Map, 0)

	for _, item := range configStructures {
		if item.Type == "list" {
			found := false
			for _, created := range lists {
				if created.Name == item.Name {
					found = true
					break
				}
			}

			if !found {
				listsToCreate = append(listsToCreate, declarations.NewList(projectId, item.Name))
			}
		}
	}

	for _, item := range lists {
		found := false
		for _, configItem := range configStructures {
			if configItem.Type == "list" {
				if item.Name == configItem.Name {
					found = true
					break
				}
			}
		}

		if !found {
			listExistsInDb = append(listExistsInDb, item)
		}
	}

	for _, item := range configStructures {
		if item.Type == "map" {
			found := false
			for _, created := range maps {
				if created.Name == item.Name {
					found = true
					break
				}
			}

			if !found {
				mapsToCreate = append(mapsToCreate, declarations.NewMap(projectId, item.Name))
			}
		}
	}

	for _, item := range maps {
		found := false
		for _, configItem := range configStructures {
			if configItem.Type == "map" {
				if item.Name == configItem.Name {
					found = true
					break
				}
			}
		}

		if !found {
			mapExistsInDb = append(mapExistsInDb, item)
		}
	}

	return listsToCreate, mapsToCreate, listExistsInDb, mapExistsInDb
}
