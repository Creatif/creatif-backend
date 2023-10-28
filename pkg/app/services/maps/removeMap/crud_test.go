package removeMap

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = ginkgo.Describe("Declaration (DELETE) a map tests", func() {
	ginkgo.It("should delete a map together with all map entries", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		handler := New(NewModel(projectId, "eng", "mapName"), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("id = ?", view.ID).First(&declarations.Map{})
		gomega.Expect(res.Error).ShouldNot(gomega.BeNil())
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())

		var mapItems []declarations.MapVariable
		res = storage.Gorm().Where("map_id = ?", view.ID).Find(&mapItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(mapItems)).Should(gomega.Equal(0))
	})
})
