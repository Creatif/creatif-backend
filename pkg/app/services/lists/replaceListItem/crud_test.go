package replaceListItem

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = ginkgo.Describe("Declaration list replace tests", func() {
	ginkgo.It("should replace a list item", func() {
		projectId := testCreateProject("project")
		variables := testCreateListAndReturnVariables(projectId, "list", 10)

		item := variables[4]
		handler := New(NewModel(projectId, "list", item["name"], Variable{
			Name:      "newName",
			Metadata:  nil,
			Groups:    nil,
			Behaviour: "readonly",
			Value:     nil,
		}))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("newName"))
		gomega.Expect(item["index"]).Should(gomega.Equal(view.Index))

		var listVariable declarations.ListVariable
		res := storage.Gorm().Where("list_id = ? AND id = ?", item["listId"], item["id"]).First(&listVariable)
		gomega.Expect(res.Error).ShouldNot(gomega.BeNil())
		gomega.Expect(res.Error).Should(gomega.MatchError(gorm.ErrRecordNotFound))
	})
})