package updateListItemByIndex

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (UPDATE) variable tests", func() {
	ginkgo.It("should update the name of the list item variable", func() {
		projectId := testCreateProject("project")
		testCreateList(projectId, "name", 100)

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, "eng", []string{"name", "behaviour"}, "name", 2, "newName", "readonly", []string{}, []byte{}, v))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Locale).Should(gomega.Equal("eng"))

		var checkModel declarations.ListVariable
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
	})

	ginkgo.It("should update the groups of a list item variable", func() {
		projectId := testCreateProject("project")
		testCreateList(projectId, "name", 100)

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(projectId, "eng", []string{"name", "groups", "value"}, "name", 6, "newName", "readonly", []string{"first", "second", "third"}, []byte{}, v))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(updated.Locale).Should(gomega.Equal("eng"))

		var checkModel declarations.ListVariable
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
	})

	ginkgo.It("should update the behaviour of a list item variable", func() {
		projectId := testCreateProject("project")
		testCreateList(projectId, "name", 100)

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(projectId, "eng", []string{"name", "behaviour", "groups"}, "name", 56, "newName", "readonly", []string{"first", "second", "third"}, []byte{}, v))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(updated.Behaviour).Should(gomega.Equal("readonly"))
		gomega.Expect(updated.Locale).Should(gomega.Equal("eng"))

		var checkModel declarations.ListVariable
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
	})
})
