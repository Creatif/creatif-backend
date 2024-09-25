package toggleProduction

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Publishing", func() {
	ginkgo.It("Should make a version production (live)", func() {
		projectId := testCreateProject("project")
		_, version := publishFullProject(projectId)

		handler := New(NewModel(projectId, version.ID), auth.NewTestingAuthentication(false, ""))
		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		var listsCount int64
		res := storage.Gorm().Raw("SELECT count(*) FROM published.published_lists").Scan(&listsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(listsCount).Should(gomega.Equal(int64(202)))

		var mapsCount int64
		res = storage.Gorm().Raw("SELECT count(*) FROM published.published_maps").Scan(&mapsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(mapsCount).Should(gomega.Equal(int64(2)))

		var referenceCount int64
		res = storage.Gorm().Raw("SELECT count(*) FROM published.published_references").Scan(&referenceCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(referenceCount).Should(gomega.Equal(int64(800)))
	})
})
