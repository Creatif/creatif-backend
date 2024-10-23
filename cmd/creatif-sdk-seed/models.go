package main

type accountVariable struct {
	name      string
	locale    string
	behaviour string
	groups    []string
	metadata  string
	value     string
}

type project struct {
	name string
	id   string
}

type account struct {
	name       string
	references []map[string]string
	imagePaths []string
	variable   accountVariable
}

func newAccountVariable(name, locale, behaviour, metadata, value string, groups []string) accountVariable {
	return accountVariable{
		name:      name,
		locale:    locale,
		behaviour: behaviour,
		groups:    groups,
		metadata:  metadata,
		value:     value,
	}
}

func newAccount(name string, references []map[string]string, imagePaths []string, variable accountVariable) account {
	return account{
		name:       name,
		references: references,
		imagePaths: imagePaths,
		variable:   variable,
	}
}

func newProject(id, name string) project {
	return project{
		name: name,
		id:   id,
	}
}
