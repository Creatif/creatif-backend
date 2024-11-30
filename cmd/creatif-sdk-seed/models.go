package main

type accountVariable struct {
	name      string
	locale    string
	behaviour string
	groups    []string
	metadata  string
	value     string
}

type propertyVariable struct {
	name      string
	locale    string
	behaviour string
	groups    []string
	metadata  string
	value     string
}

type account struct {
	name        string
	connections []map[string]string
	imagePaths  []string
	variable    accountVariable
}

type property struct {
	name        string
	connections []map[string]string
	imagePaths  []string
	variable    propertyVariable
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

func newAccount(name string, connections []map[string]string, imagePaths []string, variable accountVariable) account {
	return account{
		name:        name,
		connections: connections,
		imagePaths:  imagePaths,
		variable:    variable,
	}
}

func newProperty(name string, connections []map[string]string, imagePaths []string, variable propertyVariable) property {
	return property{
		name:        name,
		connections: connections,
		imagePaths:  imagePaths,
		variable:    variable,
	}
}

func newPropertyVariable(name, locale, behaviour, metadata, value string, groups []string) propertyVariable {
	return propertyVariable{
		name:      name,
		locale:    locale,
		behaviour: behaviour,
		groups:    groups,
		metadata:  metadata,
		value:     value,
	}
}
