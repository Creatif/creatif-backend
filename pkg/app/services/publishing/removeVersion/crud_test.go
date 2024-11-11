package removeVersion

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Publishing", func() {
	ginkgo.It("Should remove a version", func() {
		projectId := testCreateProject("project")
		_, version := publishFullProject(projectId)

		handler := New(NewModel(projectId, version.ID), auth.NewTestingAuthentication(false, ""))
		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		var listsCount int64
		res := storage.Gorm().Raw("SELECT count(*) FROM published.published_lists").Scan(&listsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(listsCount).Should(gomega.Equal(int64(0)))

		var mapsCount int64
		res = storage.Gorm().Raw("SELECT count(*) FROM published.published_maps").Scan(&mapsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(mapsCount).Should(gomega.Equal(int64(0)))
	})
})
