package paginateListItems

import "creatif/pkg/lib/sdk"

func placeGroupsInItems(items []Item) ([]Item, error) {
	normalizedGroups, err := getGroups(sdk.Map(items, func(idx int, value Item) string {
		return value.ItemID
	}))
	if err != nil {
		return nil, err
	}

	for i, item := range items {
		if _, ok := normalizedGroups[item.ItemID]; ok {
			item.Groups = normalizedGroups[item.ItemID]
			items[i] = item
		}
	}

	return items, nil
}
