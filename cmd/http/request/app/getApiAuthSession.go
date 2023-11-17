package app

import "github.com/microcosm-cc/bluemonday"

type GetApiAuthSession struct {
	Session string `param:"session"`
}

func SanitizeGetApiAuthSession(model GetApiAuthSession) GetApiAuthSession {
	p := bluemonday.StrictPolicy()
	model.Session = p.Sanitize(model.Session)

	return model
}
