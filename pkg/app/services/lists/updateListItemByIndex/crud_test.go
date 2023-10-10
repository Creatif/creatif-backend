package updateListItemByIndex

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (UPDATE) variable tests", func() {
	ginkgo.It("should update the name of the list item variable", func() {
		projectId := testCreateProject("project")
		testCreateList(projectId, "name", 100)

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.index AS index FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.name = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			"name",
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, []string{"name", "behaviour"}, "name", singleItem.Index, "newName", "readonly", []string{}, []byte{}, v))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))

		var checkModel declarations.ListVariable
		res = storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
	})

	ginkgo.It("should update the groups of a list item variable", func() {
		projectId := testCreateProject("project")
		testCreateList(projectId, "name", 100)

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.index AS index FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.name = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			"name",
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(projectId, []string{"name", "groups", "value"}, "name", singleItem.Index, "newName", "readonly", []string{"first", "second", "third"}, []byte{}, v))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))

		var checkModel declarations.ListVariable
		res = storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
	})

	ginkgo.It("should update the behaviour of a list item variable", func() {
		projectId := testCreateProject("project")
		testCreateList(projectId, "name", 100)

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.index AS index FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.name = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			"name",
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(projectId, []string{"name", "behaviour", "groups"}, "name", singleItem.Index, "newName", "readonly", []string{"first", "second", "third"}, []byte{}, v))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(updated.Behaviour).Should(gomega.Equal("readonly"))

		var checkModel declarations.ListVariable
		res = storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
	})
})
