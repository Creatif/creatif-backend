package publicApi

import (
	"strings"
)

func resolveListOptions(options string) GetListItemByIDOptions {
	opts := GetListItemByIDOptions{}

	splitOptions := strings.Split(options, ",")
	for _, opt := range splitOptions {
		keyValueOption := strings.Split(opt, ":")
		// invalid option, skip
		if len(keyValueOption) != 2 {
			continue
		}

		key := keyValueOption[0]
		value := keyValueOption[1]

		if key == "valueOnly" {
			if value == "true" {
				opts.ValueOnly = true
			} else {
				opts.ValueOnly = false
			}
		}
	}

	return opts
}

func resolveMapOptions(options string) GetMapItemByIDOptions {
	opts := GetMapItemByIDOptions{}

	splitOptions := strings.Split(options, ",")
	for _, opt := range splitOptions {
		keyValueOption := strings.Split(opt, ":")
		// invalid option, skip
		if len(keyValueOption) != 2 {
			continue
		}

		key := keyValueOption[0]
		value := keyValueOption[1]

		if key == "valueOnly" {
			if value == "true" {
				opts.ValueOnly = true
			} else {
				opts.ValueOnly = false
			}
		}
	}

	return opts
}
