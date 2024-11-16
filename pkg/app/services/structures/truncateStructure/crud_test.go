package truncateStructure

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Truncating structures", func() {
	ginkgo.It("should truncate map structure", ginkgo.Label("destructive"), func() {
		p := testCreateProject("name")
		m := testCreateMap(p, "map")

		for i := 0; i < 100; i++ {
			testAddToMap(p, m.ID, fmt.Sprintf("map-%d", i), []connections.Connection{}, []string{})
		}

		handler := New(NewModel(p, m.ID, "map"), auth.NewTestingAuthentication(false, ""))

		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		var count int
		res := storage.Gorm().
			Raw(fmt.Sprintf(`
SELECT COUNT(mv.id) FROM %s AS mv
INNER JOIN %s AS m ON project_id = ? AND m.id = ? AND mv.map_id = m.id
`,
				(declarations.MapVariable{}).TableName(),
				(declarations.Map{}).TableName(),
			), p, m.ID).Scan(&count)

		var refCount int
		res = storage.Gorm().
			Raw(fmt.Sprintf(`
SELECT COUNT(id) FROM %s WHERE project_id = ? AND (parent_structure_id = ? OR child_structure_id = ?)
`,
				(declarations.Reference{}).TableName(),
			), p, m.ID, m.ID).Scan(&count)

		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(refCount).Should(gomega.Equal(0))
	})

	ginkgo.It("should truncate list structure", ginkgo.Label("destructive"), func() {
		p := testCreateProject("name")
		m := testCreateList(p, "list")

		for i := 0; i < 100; i++ {
			testAddToList(p, m.ID, fmt.Sprintf("list-%d", i), []connections.Connection{}, []string{})
		}

		handler := New(NewModel(p, m.ID, "list"), auth.NewTestingAuthentication(false, ""))

		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		var count int
		res := storage.Gorm().
			Raw(fmt.Sprintf(`
SELECT COUNT(mv.id) FROM %s AS mv
INNER JOIN %s AS m ON project_id = ? AND m.id = ? AND mv.list_id = m.id
`,
				(declarations.ListVariable{}).TableName(),
				(declarations.List{}).TableName(),
			), p, m.ID).Scan(&count)

		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(count).Should(gomega.Equal(0))

		var refCount int
		res = storage.Gorm().
			Raw(fmt.Sprintf(`
SELECT COUNT(id) FROM %s WHERE project_id = ? AND (parent_structure_id = ? OR child_structure_id = ?)
`,
				(declarations.Reference{}).TableName(),
			), p, m.ID, m.ID).Scan(&count)

		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(refCount).Should(gomega.Equal(0))
	})
})
