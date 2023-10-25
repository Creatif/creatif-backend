package updateMapVariable

import (
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/datatypes"
)

var _ = ginkgo.Describe("Declaration (UPDATE) map entry tests", func() {
	ginkgo.It("should update an entry in the map by replacing it completely", func() {
		projectId := testCreateProject("project")
		m := testCreateMap(projectId, "map", 10)

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, "eng", m.Name, "name-0", []string{"metadata", "groups", "behaviour", "value"}, VariableModel{
			Name:      "name-0",
			Metadata:  b,
			Groups:    []string{"updated1", "updated2", "updated3"},
			Behaviour: "readonly",
			Value:     v,
		}))

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal(m.Name))
		gomega.Expect(view.ID).Should(gomega.Equal(m.ID))
		gomega.Expect(view.ProjectID).Should(gomega.Equal(m.ProjectID))

		var metadata string
		gomega.Expect(json.Unmarshal(view.Variable.Metadata.(datatypes.JSON), &metadata)).Should(gomega.BeNil())

		var value string
		gomega.Expect(json.Unmarshal(view.Variable.Value.(datatypes.JSON), &value)).Should(gomega.BeNil())

		entry := view.Variable
		gomega.Expect(entry.Name).Should(gomega.Equal("name-0"))
		gomega.Expect(metadata).Should(gomega.Equal("this is metadata"))
		gomega.Expect(value).Should(gomega.Equal("this is value"))
		gomega.Expect(entry.Behaviour).Should(gomega.Equal("readonly"))
		gomega.Expect(sdk.Includes(entry.Groups, "updated1")).Should(gomega.Equal(true))
		gomega.Expect(sdk.Includes(entry.Groups, "updated2")).Should(gomega.Equal(true))
		gomega.Expect(sdk.Includes(entry.Groups, "updated3")).Should(gomega.Equal(true))
	})

	ginkgo.It("should fail updating a map variable because of invalid number of groups", func() {
		projectId := testCreateProject("project")
		m := testCreateMap(projectId, "map", 10)

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, "eng", m.Name, "name-0", []string{"metadata", "groups", "behaviour", "value"}, VariableModel{
			Name:      "name-0",
			Metadata:  b,
			Groups:    []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18"},
			Behaviour: "readonly",
			Value:     v,
		}))

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["groups"]).ShouldNot(gomega.BeEmpty())
	})
})
