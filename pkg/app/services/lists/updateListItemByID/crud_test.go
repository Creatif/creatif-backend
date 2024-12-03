package updateListItemByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/appErrors"
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
			[]connections.Connection{},
			[]string{},
		), auth.NewTestingAuthentication(false, ""))

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
		handler := New(NewModel(projectId, "eng", []string{"name", "groups", "value"}, view.ID, singleItem.ID, "newName", "readonly", []string{g[0], g[1], g[2]}, []byte{}, v, []connections.Connection{}, []string{}), auth.NewTestingAuthentication(false, ""))

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
			[]connections.Connection{},
			[]string{},
		), auth.NewTestingAuthentication(false, ""))

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
			[]connections.Connection{},
			[]string{},
		),
			auth.NewTestingAuthentication(false, ""),
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
			[]connections.Connection{},
			[]string{},
		),
			auth.NewTestingAuthentication(false, ""),
		)

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["behaviourReadonly"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should update variable with mixed connections", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId, []string{"one", "two", "three"})
		view := testCreateList(projectId, "name", 100, false, "modifiable")

		connectionList := testCreateList(projectId, "connection list", 0, false, "modifiable")
		connectionMap := testCreateMap(projectId, "connection map")

		createMapConnections := func() []connections.Connection {
			connsViews := make([]addToMap.View, 5)
			for i := 0; i < 5; i++ {
				connsViews[i] = testAddToMap(projectId, fmt.Sprintf("variable-%d", i), connectionMap.ID, []connections.Connection{}, []string{})
			}

			conns := make([]connections.Connection, 5)
			for i, c := range connsViews {
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-map-%d", i),
					StructureType: "map",
					VariableID:    c.ID,
				}
			}

			return conns
		}

		createListConnections := func() []connections.Connection {
			connsViews := make([]addToList.View, 5)
			for i := 0; i < 5; i++ {
				connsViews[i] = testAddToList(projectId, fmt.Sprintf("variable-%d", i), connectionList.ID, []connections.Connection{}, []string{})
			}

			conns := make([]connections.Connection, 5)
			for i, c := range connsViews {
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-list-%d", i),
					StructureType: "list",
					VariableID:    c.ID,
				}
			}

			return conns
		}

		var conns []connections.Connection
		conns = append(conns, createListConnections()...)
		conns = append(conns, createMapConnections()...)

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
			[]string{"name", "behaviour", "connections"},
			view.ID,
			singleItem.ID,
			"newName",
			"readonly",
			[]string{},
			[]byte{},
			v,
			conns,
			[]string{},
		), auth.NewTestingAuthentication(false, ""))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		var connectionsCount int
		res = storage.Gorm().Raw(
			fmt.Sprintf("SELECT COUNT(*) as count FROM %s WHERE parent_variable_id = ?", (declarations.Connection{}).TableName()),
			updated.ID,
		).Scan(&connectionsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(connectionsCount).Should(gomega.Equal(10))
	})

	ginkgo.It("should update variable with mixed connections when the connections is empty", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId, []string{"one", "two", "three"})
		view := testCreateList(projectId, "name", 100, false, "modifiable")

		connectionList := testCreateList(projectId, "connection list", 0, false, "modifiable")
		connectionMap := testCreateMap(projectId, "connection map")

		createMapConnections := func() []connections.Connection {
			connsViews := make([]addToMap.View, 5)
			for i := 0; i < 5; i++ {
				connsViews[i] = testAddToMap(projectId, fmt.Sprintf("variable-%d", i), connectionMap.ID, []connections.Connection{}, []string{})
			}

			conns := make([]connections.Connection, 5)
			for i, c := range connsViews {
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-map-%d", i),
					StructureType: "map",
					VariableID:    c.ID,
				}
			}

			return conns
		}

		createListConnections := func() []connections.Connection {
			connsViews := make([]addToList.View, 5)
			for i := 0; i < 5; i++ {
				connsViews[i] = testAddToList(projectId, fmt.Sprintf("variable-%d", i), connectionList.ID, []connections.Connection{}, []string{})
			}

			conns := make([]connections.Connection, 5)
			for i, c := range connsViews {
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-list-%d", i),
					StructureType: "list",
					VariableID:    c.ID,
				}
			}

			return conns
		}

		var conns []connections.Connection
		conns = append(conns, createListConnections()...)
		conns = append(conns, createMapConnections()...)

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
			[]string{"name", "connections"},
			view.ID,
			singleItem.ID,
			"newName",
			"modifiable",
			[]string{},
			[]byte{},
			v,
			conns,
			[]string{},
		), auth.NewTestingAuthentication(false, ""))

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		var connectionsCount int
		res = storage.Gorm().Raw(
			fmt.Sprintf("SELECT COUNT(*) as count FROM %s WHERE parent_variable_id = ?", (declarations.Connection{}).TableName()),
			updated.ID,
		).Scan(&connectionsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(connectionsCount).Should(gomega.Equal(10))

		handler = New(NewModel(
			projectId,
			"eng",
			[]string{"name", "connections"},
			view.ID,
			singleItem.ID,
			"newName",
			"modifiable",
			[]string{},
			[]byte{},
			v,
			[]connections.Connection{},
			[]string{},
		), auth.NewTestingAuthentication(false, ""))

		updated, err = handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		res = storage.Gorm().Raw(
			fmt.Sprintf("SELECT COUNT(*) as count FROM %s WHERE parent_variable_id = ?", (declarations.Connection{}).TableName()),
			updated.ID,
		).Scan(&connectionsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(connectionsCount).Should(gomega.Equal(0))
	})
})
