package create

import (
	"encoding/json"
)

func isOrdinaryVisit(existing map[string]string, incoming map[string]string) bool {
	isGroupVisit := func(existing map[string]string, incoming map[string]string) bool {
		if existing["type"] == "visit" &&
			existing["subType"] == "groups" &&
			incoming["type"] == "visit" &&
			incoming["subType"] == "groups" {
			return true
		}

		return false
	}

	isApiVisit := func(existing map[string]string, incoming map[string]string) bool {
		if existing["type"] == "visit" &&
			existing["subType"] == "api" &&
			incoming["type"] == "visit" &&
			incoming["subType"] == "api" {
			return true
		}

		return false
	}

	isMapStructureVisit := func(existing map[string]string, incoming map[string]string) bool {
		if existing["type"] == "visit" &&
			existing["subType"] == "mapStructures" &&
			incoming["type"] == "visit" &&
			incoming["subType"] == "mapStructures" {
			return true
		}

		return false
	}

	isListStructureVisit := func(existing map[string]string, incoming map[string]string) bool {
		if existing["type"] == "visit" &&
			existing["subType"] == "listStructures" &&
			incoming["type"] == "visit" &&
			incoming["subType"] == "listStructures" {
			return true
		}

		return false
	}

	if isGroupVisit(existing, incoming) {
		return true
	}

	if isApiVisit(existing, incoming) {
		return true
	}

	if isMapStructureVisit(existing, incoming) {
		return true
	}

	if isListStructureVisit(existing, incoming) {
		return true
	}

	return false
}

func decideToCreateNewActivity(allQueries []DataQuery, incoming []byte) (bool, error) {
	for _, query := range allQueries {
		var lastWrittenDataQuery map[string]string
		var incomingDataQuery map[string]string
		if err := json.Unmarshal(query.Data, &lastWrittenDataQuery); err != nil {
			return false, err
		}

		if err := json.Unmarshal(incoming, &incomingDataQuery); err != nil {
			return false, err
		}

		if isOrdinaryVisit(lastWrittenDataQuery, incomingDataQuery) {
			return false, nil
		}
	}

	return true, nil
}
