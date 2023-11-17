package removeMapEntry

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = ginkgo.Describe("Declaration (DELETE) map entry tests", func() {
	ginkgo.It("should delete a map entry by name", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		variables := view.Variables
		entryName := variables[0]["name"]
		handler := New(NewModel(projectId, "eng", "mapName", "", "", entryName, "", ""), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("map_id = ? AND name = ?", view.ID, entryName).First(&declarations.MapVariable{})
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())
	})

	ginkgo.It("should delete a map entry multiple fields (1)", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		variables := view.Variables
		shortId := variables[0]["shortID"]
		handler := New(NewModel(projectId, "eng", "", view.ID, "", "", "", shortId), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("map_id = ? AND short_id = ?", view.ID, shortId).First(&declarations.MapVariable{})
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())
	})

	ginkgo.It("should delete a map entry multiple fields (2)", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		variables := view.Variables
		variableName := variables[0]["name"]
		handler := New(NewModel(projectId, "eng", "", "", view.ShortID, variableName, "", ""), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("map_id = ? AND name = ?", view.ID, variableName).First(&declarations.MapVariable{})
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())
	})

	ginkgo.It("should delete a map entry multiple fields (3)", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		variables := view.Variables
		shortId := variables[0]["shortID"]
		handler := New(NewModel(projectId, "eng", view.Name, "", "", "", "", shortId), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("map_id = ? AND short_id = ?", view.ID, shortId).First(&declarations.MapVariable{})
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())
	})
})
