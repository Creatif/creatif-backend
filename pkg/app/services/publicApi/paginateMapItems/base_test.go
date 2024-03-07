package paginateMapItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/lists/addToList"
	createList2 "creatif/pkg/app/services/lists/createList"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/app/services/maps/mapCreate"
	createProject2 "creatif/pkg/app/services/projects/createProject"
	"creatif/pkg/app/services/publishing/publish"
	"creatif/pkg/app/services/publishing/toggleProduction"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	storage2 "creatif/pkg/lib/storage"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

func loadEnv() {
	err := godotenv.Load("../../../../../.env")

	if err != nil {
		log.Fatal(err)
	}
}

var GomegaRegisterFailHandler = gomega.RegisterFailHandler
var GinkgoFail = ginkgo.Fail
var GinkgoRunSpecs = ginkgo.RunSpecs
var GinkgoAfterHandler = ginkgo.AfterEach
var GinkgoAfterSuite = ginkgo.AfterSuite

func TestApi(t *testing.T) {
	GomegaRegisterFailHandler(GinkgoFail)
	GinkgoRunSpecs(t, "Publishing tests")
}

func runLogger() {
	if err := logger.BuildLoggers(os.Getenv("LOG_DIRECTORY")); err != nil {
		log.Fatalln(fmt.Sprintf("Cannot createProject logger: %s", err.Error()))
	}

	logger.Info("Health info logger health check... Ignore!")
	logger.Warn("Health warning logger health check... Ignore!")
	logger.Error("Health error logger health check... Ignore!")
}

var _ = ginkgo.BeforeSuite(func() {
	loadEnv()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Zagreb",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PORT"),
	)

	err := storage2.Connect(dsn)

	if err != nil {
		log.Fatalln(err)
	}

	gomega.Expect(locales.Store()).Should(gomega.BeNil())
})

var _ = GinkgoAfterSuite(func() {
	sql, err := storage2.SQLDB()
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("Could not retreive DB instance: %s", err))
	}

	if err := sql.Close(); err != nil {
		ginkgo.Fail(fmt.Sprintf("Could not close database connection: %s", err))
	}
})

var _ = GinkgoAfterHandler(func() {
	res := storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", domain.PROJECT_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.MAP_VARIABLES))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.VARIABLE_MAP))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.LIST_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.LIST_VARIABLES_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", domain.USERS_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.REFERENCE_TABLES))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.GROUPS_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.VARIABLE_GROUPS_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", domain.PUBLISHED_LISTS_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", domain.PUBLISHED_MAPS_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", domain.VERSION_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", domain.PUBLISHED_REFERENCES_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
})

func testAssertErrNil(err error) {
	gomega.Expect(err).Should(gomega.BeNil())
}

func testAssertIDValid(id string) {
	gomega.Expect(id).ShouldNot(gomega.BeEmpty())
	_, err := ulid.Parse(id)
	gomega.Expect(err).Should(gomega.BeNil())
}

func testCreateProject(name string) string {
	handler := createProject2.New(createProject2.NewModel(name), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

	model, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(model.ID)

	gomega.Expect(model.Name).Should(gomega.Equal(name))

	return model.ID
}

func testCreateMap(projectId, name string) mapCreate.View {
	entries := make([]mapCreate.VariableModel, 0)

	handler := mapCreate.New(mapCreate.NewModel(projectId, name, entries), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID)

	gomega.Expect(name).Should(gomega.Equal(view.Name))

	return view
}

func testCreateList(projectId, name string) createList2.View {
	handler := createList2.New(createList2.NewModel(projectId, name, []createList2.Variable{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

	list, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(list.ID)

	gomega.Expect(list.Name).Should(gomega.Equal(name))

	return list
}

func testCreateGroups(projectId string, numOfGroups int) []addGroups.View {
	groups := make([]addGroups.GroupModel, numOfGroups)
	for i := 0; i < numOfGroups; i++ {
		groups[i] = addGroups.GroupModel{
			ID:     "",
			Name:   fmt.Sprintf("group-%d", i),
			Type:   "new",
			Action: "create",
		}
	}

	l := logger.NewLogBuilder()

	handler := addGroups.New(addGroups.NewModel(projectId, groups), auth.NewTestingAuthentication(false, projectId), l)
	model, err := handler.Handle()
	testAssertErrNil(err)

	return model
}

func testAddToMap(projectId, name, variableName string, references []shared.Reference, groups []string) addToMap.View {
	variableModel := addToMap.VariableModel{
		Name:      variableName,
		Metadata:  nil,
		Groups:    groups,
		Value:     nil,
		Locale:    "eng",
		Behaviour: "modifiable",
	}

	model := addToMap.NewModel(projectId, name, variableModel, references)
	handler := addToMap.New(model, auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

	view, err := handler.Handle()
	gomega.Expect(err).Should(gomega.BeNil())

	return view
}

func testAddToList(projectId, name, variableName string, references []shared.Reference, groups []string) addToList.View {
	variableModel := addToList.VariableModel{
		Name:      variableName,
		Metadata:  nil,
		Groups:    groups,
		Value:     nil,
		Locale:    "eng",
		Behaviour: "modifiable",
	}

	model := addToList.NewModel(projectId, name, variableModel, references)
	handler := addToList.New(model, auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

	view, err := handler.Handle()
	gomega.Expect(err).Should(gomega.BeNil())

	return view
}

func publishFullProject(projectId string) ([]addToMap.View, publish.View) {
	groups := testCreateGroups(projectId, 5)

	paginationList := testCreateMap(projectId, "paginationMap")

	referenceList := testCreateList(projectId, "referenceList")
	referenceListItem1 := testAddToList(projectId, referenceList.ID, "reference-list-1", []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
		return value.ID
	}))
	referenceListItem2 := testAddToList(projectId, referenceList.ID, "reference-map-2", []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
		return value.ID
	}))

	referenceMap := testCreateMap(projectId, "referenceMap")
	referenceMapItem1 := testAddToMap(projectId, referenceMap.ID, "reference-map-1", []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
		return value.ID
	}))
	referenceMapItem2 := testAddToMap(projectId, referenceMap.ID, "reference-map-2", []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
		return value.ID
	}))

	views := make([]addToMap.View, 0)
	for i := 0; i < 200; i++ {
		addToListModel := testAddToMap(projectId, paginationList.ID, fmt.Sprintf("name-%d", i), []shared.Reference{
			{
				Name:          "first",
				StructureName: referenceMap.Name,
				StructureType: "map",
				VariableID:    referenceMapItem1.ID,
			},
			{
				Name:          "second",
				StructureName: referenceMap.Name,
				StructureType: "map",
				VariableID:    referenceMapItem2.ID,
			},
			{
				Name:          "third",
				StructureName: paginationList.Name,
				StructureType: "list",
				VariableID:    referenceListItem1.ID,
			},
			{
				Name:          "fourth",
				StructureName: paginationList.Name,
				StructureType: "list",
				VariableID:    referenceListItem2.ID,
			},
		}, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		views = append(views, addToListModel)
	}

	handler := publish.New(publish.NewModel(projectId, "v1"), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
	model, err := handler.Handle()
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(model.ID).ShouldNot(gomega.BeEmpty())
	gomega.Expect(model.Name).ShouldNot(gomega.BeEmpty())

	toggleHandler := toggleProduction.New(toggleProduction.NewModel(projectId, model.ID), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
	_, err = toggleHandler.Handle()
	gomega.Expect(err).Should(gomega.BeNil())

	return views, model
}
