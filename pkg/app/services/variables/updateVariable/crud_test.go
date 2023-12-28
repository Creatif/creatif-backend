package updateVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (UPDATE) variable tests", func() {
	ginkgo.It("should update the name of the declaration variable", func() {
		projectId := testCreateProject("project")
		view := testCreateBasicDeclarationTextVariable(projectId, "name", "modifiable")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, []string{"name", "behaviour"}, view.ShortID, "newName", "readonly", []string{}, []byte{}, v, "eng"), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(updated.ID))
		gomega.Expect(view.ProjectID).Should(gomega.Equal(projectId))
		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Locale).Should(gomega.Equal("eng"))

		var checkModel declarations.Variable
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
		gomega.Expect(checkModel.Value).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should update the groups of the declaration variable", func() {
		projectId := testCreateProject("project")
		view := testCreateBasicDeclarationTextVariable(projectId, "name", "modifiable")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(projectId, []string{"name", "groups", "value"}, view.ID, "newName", "readonly", []string{"first", "second", "third"}, []byte{}, v, "eng"), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(updated.ID))
		gomega.Expect(view.ProjectID).Should(gomega.Equal(projectId))
		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Locale).Should(gomega.Equal("eng"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))

		var checkModel declarations.Variable
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
	})

	ginkgo.It("should update the behaviour of the declaration variable", func() {
		projectId := testCreateProject("project")
		view := testCreateBasicDeclarationTextVariable(projectId, "name", "modifiable")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(projectId, []string{"name", "behaviour", "groups"}, view.ShortID, "newName", "readonly", []string{"first", "second", "third"}, []byte{}, v, "eng"), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		updated, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(updated.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(updated.ID))
		gomega.Expect(view.ProjectID).Should(gomega.Equal(projectId))
		gomega.Expect(updated.Name).Should(gomega.Equal("newName"))
		gomega.Expect(updated.Locale).Should(gomega.Equal("eng"))
		gomega.Expect(updated.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(updated.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(updated.Behaviour).Should(gomega.Equal("readonly"))

		var checkModel declarations.Variable
		res := storage.Gorm().Table(checkModel.TableName()).Where("id = ?", updated.ID).First(&checkModel)
		testAssertErrNil(res.Error)

		gomega.Expect(checkModel.Name).Should(gomega.Equal("newName"))
		gomega.Expect(checkModel.Groups).Should(gomega.HaveLen(3))
		gomega.Expect(checkModel.Groups[0]).Should(gomega.Equal("first"))
		gomega.Expect(checkModel.Behaviour).Should(gomega.Equal("readonly"))
	})

	ginkgo.It("should fail updating groups if the total number of groups is > 20", func() {
		projectId := testCreateProject("project")
		view := testCreateBasicDeclarationTextVariable(projectId, "name", "modifiable")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(
			projectId,
			[]string{"name", "behaviour", "groups"},
			view.ID,
			"newName",
			"readonly",
			[]string{"1", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19"},
			[]byte{}, v, "eng"),
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

	ginkgo.It("should fail updating readonly variable", func() {
		projectId := testCreateProject("project")
		view := testCreateBasicDeclarationTextVariable(projectId, "name", "readonly")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(
			projectId,
			[]string{"name", "behaviour", "groups"},
			view.ShortID,
			"newName",
			"readonly",
			[]string{"1", "1", "2"},
			[]byte{}, v, "eng"),
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

	ginkgo.It("should fail update if locale does not exist", func() {
		projectId := testCreateProject("project")
		view := testCreateBasicDeclarationTextVariable(projectId, "name", "readonly")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(
			projectId,
			[]string{"name", "behaviour", "groups"},
			view.ID,
			"newName",
			"readonly",
			[]string{"1", "1", "2"},
			[]byte{},
			v,
			"",
		),
			auth.NewTestingAuthentication(false),
			logger.NewLogBuilder(),
		)

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["nameLocaleValid"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should fail updating name if name does not exist", func() {
		projectId := testCreateProject("project")
		view := testCreateBasicDeclarationTextVariable(projectId, "name", "readonly")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(
			projectId,
			[]string{"name", "behaviour", "groups"},
			view.ShortID,
			"",
			"readonly",
			[]string{"1", "1", "2"},
			[]byte{},
			v,
			"eng",
		),
			auth.NewTestingAuthentication(false),
			logger.NewLogBuilder(),
		)

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["nameLocaleValid"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should fail variable update id does not exist", func() {
		projectId := testCreateProject("project")
		testCreateBasicDeclarationTextVariable(projectId, "name", "readonly")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(
			projectId,
			[]string{"name", "behaviour", "groups"},
			"not exists",
			"name",
			"readonly",
			[]string{"1", "1", "2"},
			[]byte{},
			v,
			"eng",
		),
			auth.NewTestingAuthentication(false),
			logger.NewLogBuilder(),
		)

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["notExists"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should fail variable with name and locale already exists", func() {
		projectId := testCreateProject("project")
		view := testCreateBasicDeclarationTextVariable(projectId, "name", "modifiable")

		m := "text value"
		v, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())
		handler := New(NewModel(
			projectId,
			[]string{"name", "behaviour", "groups"},
			view.ShortID,
			view.Name,
			"modifiable",
			[]string{"1", "1", "2"},
			[]byte{},
			v,
			"eng",
		),
			auth.NewTestingAuthentication(false),
			logger.NewLogBuilder(),
		)

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["exists"]).ShouldNot(gomega.BeEmpty())
	})
})
