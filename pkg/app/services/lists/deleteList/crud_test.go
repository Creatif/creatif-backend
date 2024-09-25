package deleteList

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list delete tests", func() {
	ginkgo.It("should delete a list by name", func() {
		projectId := testCreateProject("project")
		listName, listId, _ := testCreateListAndReturnNameAndID(projectId, "name", 100)

		handler := New(NewModel(projectId, listName, "", ""), auth.NewTestingAuthentication(false, ""))
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var listItems []declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listId).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(listItems)).Should(gomega.Equal(0))
	})

	ginkgo.It("should delete a list by id", func() {
		projectId := testCreateProject("project")
		_, listId, _ := testCreateListAndReturnNameAndID(projectId, "name", 100)

		handler := New(NewModel(projectId, "", listId, ""), auth.NewTestingAuthentication(false, ""))
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var listItems []declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listId).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(listItems)).Should(gomega.Equal(0))
	})

	ginkgo.It("should delete a list by shortID", func() {
		projectId := testCreateProject("project")
		_, listID, shortID := testCreateListAndReturnNameAndID(projectId, "name", 100)

		handler := New(NewModel(projectId, "", "", shortID), auth.NewTestingAuthentication(false, ""))
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var listItems []declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listID).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(listItems)).Should(gomega.Equal(0))
	})
})
