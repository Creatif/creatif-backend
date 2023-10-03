package createList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list create tests", func() {
	ginkgo.It("should create a list", func() {
		projectId := testCreateProject("project")
		variables := make([]Variable, 5)
		for i := 0; i < 5; i++ {
			variables[i] = Variable{
				Name:      fmt.Sprintf("one-%d", i),
				Metadata:  nil,
				Groups:    nil,
				Behaviour: "readonly",
				Value:     nil,
			}
		}

		handler := New(NewModel(projectId, "list", variables))

		list, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(list.ID)

		gomega.Expect(list.Name).Should(gomega.Equal("list"))

		var savedVariables []declarations.ListVariable
		storage.Gorm().Where("list_id = ?", list.ID).Find(&savedVariables)

		gomega.Expect(len(savedVariables)).Should(gomega.Equal(5))
		for i := 1; i <= 5; i++ {
			gomega.Expect(savedVariables[i-1].Index).Should(gomega.Equal(int64(i)))
		}
	})
})
