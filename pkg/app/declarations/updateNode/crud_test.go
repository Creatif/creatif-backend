package updateNode

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (UPDATE) node tests", func() {
	ginkgo.It("should update the name of the declaration node", func() {
		view := testCreateBasicDeclarationTextNode("name", "modifiable")

		handler := New(NewModel([]string{"name", "behaviour"}, "name", "newName", "readonly", []string{}, []byte{}))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(updated.ID))
		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))

		var checkModel declarations.Node
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
	})

	ginkgo.It("should update the groups of the declaration node", func() {
		view := testCreateBasicDeclarationTextNode("name", "modifiable")

		handler := New(NewModel([]string{"name", "groups"}, "name", "newName", "readonly", []string{"first", "second", "third"}, []byte{}))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(updated.ID))
		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))

		var checkModel declarations.Node
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
	})

	ginkgo.It("should update the behaviour of the declaration node", func() {
		view := testCreateBasicDeclarationTextNode("name", "modifiable")

		handler := New(NewModel([]string{"name", "behaviour", "groups"}, "name", "newName", "readonly", []string{"first", "second", "third"}, []byte{}))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(updated.ID))
		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(updated.Behaviour).Should(gomega.Equal("readonly"))

		var checkModel declarations.Node
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
	})
})
