package main

import (
	"fmt"
	"strings"
)

type reporter struct {
	projectIds      []string
	numOfAccounts   int
	numOfProperties int
}

func newReporter() *reporter {
	return &reporter{
		projectIds:      make([]string, 0),
		numOfAccounts:   0,
		numOfProperties: 0,
	}
}

func (r *reporter) Report() {
	fmt.Printf(`
Seeding statistics and useful information

Number of projects: %d
Number of accounts: %d
Number of properties: %d

This seed is not intended to be used in the UI project but it can be. 
If you want, just replace the project ID in the URL with any of the below project IDs

%s

`, len(r.projectIds), r.numOfAccounts, r.numOfProperties, strings.Join(r.projectIds, "\n"))
}

func (r *reporter) AddProjectID(projectId string) {
	r.projectIds = append(r.projectIds, projectId)
}

func (r *reporter) AddAccount() {
	r.numOfAccounts += 1
}

func (r *reporter) AddProperty() {
	r.numOfProperties += 1
}