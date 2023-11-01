package shared

import "fmt"

func determineIDPlaceholder(alias, fieldName, name, id, shortId string, namedPlaceholder bool) string {
	var t string
	v := "?"
	if id != "" {
		if namedPlaceholder {
			v = "@id"
		}
		t = fmt.Sprintf("id = %s", v)
	} else if name != "" {
		if namedPlaceholder {
			v = "@name"
		}
		
		if fieldName != "" {
			v = fmt.Sprintf("@%s", fieldName)
		}
		t = fmt.Sprintf("name = %s", v)
	} else if shortId != "" {
		if namedPlaceholder {
			v = "@id"
		}
		t = fmt.Sprintf("short_id = %s", v)
	}

	if alias != "" {
		t = fmt.Sprintf("%s.%s", alias, t)
	}

	return t
}

func determineIDNamedPlaceholder(alias, name, id, shortId string) string {
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

func determineIDValue(name, id, shortID string) string {
	if id != "" {
		return id
	} else if name != "" {
		return name
	}

	return shortID
}

func DetermineID(alias, name, id, shortID string) (string, string) {
	return determineIDPlaceholder(alias, "", name, id, shortID, false), determineIDValue(name, id, shortID)
}

func DetermineIDWithNamedPlaceholder(alias, fieldName, name, id, shortID string) (string, string) {
	return determineIDPlaceholder(alias, fieldName, name, id, shortID, true), determineIDValue(name, id, shortID)
}
