package shared

import "fmt"

func DetermineIDPlaceholder(alias, name, id, shortId string) string {
	var t string
	if id != "" {
		t = fmt.Sprintf("id = ?")
	} else if name != "" {
		t = fmt.Sprintf("name = ?")
	} else if shortId != "" {
		t = fmt.Sprintf("short_id = ?")
	}

	if alias != "" {
		t = fmt.Sprintf("%s.%s", alias, t)
	}

	return t
}

func DetermineIDValue(name, id, shortID string) string {
	if id != "" {
		return id
	} else if name != "" {
		return name
	}

	return shortID
}

func DetermineID(alias, name, id, shortID string) (string, string) {
	return DetermineIDPlaceholder(alias, name, id, shortID), DetermineIDValue(name, id, shortID)
}
