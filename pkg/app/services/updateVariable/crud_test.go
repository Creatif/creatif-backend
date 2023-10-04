package updateVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (UPDATE) variable tests", func() {
	ginkgo.It("should update the name of the declaration variable", func() {
		projectId := testCreateProject("project")
		view := testCreateBasicDeclarationTextVariable(projectId, "name", "modifiable")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, []string{"name", "behaviour"}, "name", "newName", "readonly", []string{}, []byte{}, v))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(updated.ID))
		gomega.Expect(view.ProjectID).Should(gomega.Equal(projectId))
		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))

		var checkModel declarations.Variable
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
		gomega.Expect(checkModel.Value).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should update the groups of the declaration variable", func() {
		projectId := testCreateProject("project")
		view := testCreateBasicDeclarationTextVariable(projectId, "name", "modifiable")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(projectId, []string{"name", "groups", "value"}, "name", "newName", "readonly", []string{"first", "second", "third"}, []byte{}, v))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(updated.ID))
		gomega.Expect(view.ProjectID).Should(gomega.Equal(projectId))
		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))

		var checkModel declarations.Variable
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
	})

	ginkgo.It("should update the behaviour of the declaration variable", func() {
		projectId := testCreateProject("project")
		view := testCreateBasicDeclarationTextVariable(projectId, "name", "modifiable")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(projectId, []string{"name", "behaviour", "groups"}, "name", "newName", "readonly", []string{"first", "second", "third"}, []byte{}, v))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(updated.ID))
		gomega.Expect(view.ProjectID).Should(gomega.Equal(projectId))
		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(updated.Behaviour).Should(gomega.Equal("readonly"))

		var checkModel declarations.Variable
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
	})
})
