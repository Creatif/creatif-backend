package updateListItemByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (UPDATE) variable tests", func() {
	ginkgo.It("should update the name of the list item variable", func() {
		projectId := testCreateProject("project")
		testCreateList(projectId, "name", 100, false, "modifiable")

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.id AS id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.name = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			"name",
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, "eng", []string{"name", "behaviour"}, "name", singleItem.ID, "newName", "readonly", []string{}, []byte{}, v), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Locale).Should(gomega.Equal("eng"))

		var checkModel declarations.ListVariable
		res = storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
	})

	ginkgo.It("should update the groups of the declaration variable", func() {
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "name", 100, false, "modifiable")

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.short_id AS short_id, lv.id AS id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.name = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			listName,
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(res.RowsAffected).ShouldNot(gomega.Equal(0))

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(projectId, "eng", []string{"name", "groups", "value"}, "name", singleItem.ShortID, "newName", "readonly", []string{"first", "second", "third"}, []byte{}, v), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Locale).Should(gomega.Equal("eng"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))

		var checkModel declarations.ListVariable
		res = storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
	})

	ginkgo.It("should update the behaviour of the declaration variable", func() {
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "name", 100, false, "modifiable")

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.id AS id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.name = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			listName,
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(projectId, "eng", []string{"name", "behaviour", "groups"}, "name", singleItem.ID, "newName", "readonly", []string{"first", "second", "third"}, []byte{}, v), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(updated.Behaviour).Should(gomega.Equal("readonly"))
		gomega.Expect(updated.Locale).Should(gomega.Equal("eng"))

		var checkModel declarations.ListVariable
		res = storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
	})

	ginkgo.It("should fail updating list variable because max number of groups", func() {
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "name", 100, true, "modifiable")

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.id AS id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.name = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			listName,
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(
			projectId,
			"eng",
			[]string{"name", "behaviour", "groups"},
			"name",
			singleItem.ID,
			"newName",
			"readonly",
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18"},
			[]byte{}, v),
			auth.NewTestingAuthentication(false),
			logger.NewLogBuilder(),
		)

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["groups"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should fail updating list variable because of invalid behaviour", func() {
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "name", 100, true, "readonly")

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.short_id, lv.id AS id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.name = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			listName,
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(
			projectId,
			"eng",
			[]string{"name", "behaviour", "groups"},
			"name",
			singleItem.ShortID,
			"newName",
			"readonly",
			[]string{"1", "2", "3", "4"},
			[]byte{}, v),
			auth.NewTestingAuthentication(false),
			logger.NewLogBuilder(),
		)

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["behaviour"]).ShouldNot(gomega.BeEmpty())
	})
})
