package updateListItemByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (UPDATE) variable tests", func() {
	ginkgo.It("should update the name of the list item variable", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId, []string{"one", "two", "three"})
		view := testCreateList(projectId, "name", 100, false, "modifiable")

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.id AS id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.id = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			view.ID,
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(
			projectId,
			"eng",
			[]string{"name", "behaviour"},
			view.ID,
			singleItem.ID,
			"newName",
			"readonly",
			[]string{},
			[]byte{},
			v,
			[]shared.UpdateReference{},
		), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

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

	ginkgo.It("should update the groups of the a list item", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, []string{"first", "second", "third", "one", "two", "three"})
		view := testCreateList(projectId, "name", 100, false, "modifiable")

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.short_id AS short_id, lv.id AS id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.short_id = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			view.ShortID,
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(res.RowsAffected).ShouldNot(gomega.Equal(0))

		g := sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		})

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(projectId, "eng", []string{"name", "groups", "value"}, view.ID, singleItem.ID, "newName", "readonly", []string{g[0], g[1], g[2]}, []byte{}, v, []shared.UpdateReference{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Locale).Should(gomega.Equal("eng"))
		gomega.Expect(len(updated.Groups)).Should(gomega.Equal(3))

		var checkModel declarations.ListVariable
		res = storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
	})

	ginkgo.It("should update the behaviour of the list item", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, []string{"one", "two", "three", "first", "second", "third"})
		view := testCreateList(projectId, "name", 100, false, "modifiable")

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.id AS id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.id = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			view.ID,
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		g := sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		})

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(
			projectId,
			"eng",
			[]string{"name", "behaviour", "groups"},
			view.ShortID,
			singleItem.ID,
			"newName",
			"readonly",
			[]string{g[0], g[1], g[2]},
			[]byte{},
			v,
			[]shared.UpdateReference{},
		), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Behaviour).Should(gomega.Equal("readonly"))
		gomega.Expect(updated.Locale).Should(gomega.Equal("eng"))

		var checkModel declarations.ListVariable
		res = storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
	})

	ginkgo.It("should fail updating list item because of non existing groups", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId, []string{"one", "two", "three"})
		view := testCreateList(projectId, "name", 100, true, "modifiable")

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.id AS id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.id = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			view.ID,
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
			view.ID,
			singleItem.ID,
			"newName",
			"readonly",
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18"},
			[]byte{},
			v,
			[]shared.UpdateReference{},
		),
			auth.NewTestingAuthentication(false, ""),
			logger.NewLogBuilder(),
		)

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["groupsExist"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should fail updating list variable because of invalid behaviour", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId, []string{"one", "two", "three"})
		view := testCreateList(projectId, "name", 100, true, "readonly")

		var singleItem declarations.ListVariable
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT lv.short_id, lv.id AS id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.id = ? AND l.project_id = ?", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()),
			view.ID,
			projectId,
		).Scan(&singleItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(
			projectId,
			"eng",
			[]string{"name", "behaviour"},
			view.ID,
			singleItem.ShortID,
			"newName",
			"readonly",
			nil,
			[]byte{},
			v,
			[]shared.UpdateReference{},
		),
			auth.NewTestingAuthentication(false, ""),
			logger.NewLogBuilder(),
		)

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["behaviourReadonly"]).ShouldNot(gomega.BeEmpty())
	})
})
