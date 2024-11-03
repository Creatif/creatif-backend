package getStructures

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
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/sdk"
	storage2 "creatif/pkg/lib/storage"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/segmentio/ksuid"
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
	GinkgoRunSpecs(t, "Public API - Get structures")
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
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", domain.ACTIVITY))
	gomega.Expect(res.Error).Should(gomega.BeNil())
})

func testAssertErrNil(err error) {
	gomega.Expect(err).Should(gomega.BeNil())
}

func testAssertIDValid(id string) {
	gomega.Expect(id).ShouldNot(gomega.BeEmpty())
	_, err := ksuid.Parse(id)
	gomega.Expect(err).Should(gomega.BeNil())
}

func testCreateProject(name string) string {
	handler := createProject2.New(createProject2.NewModel(name), auth.NewTestingAuthentication(false, ""))

	model, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(model.ID)

	gomega.Expect(model.Name).Should(gomega.Equal(name))

	return model.ID
}

func testCreateMap(projectId, name string) mapCreate.View {
	entries := make([]mapCreate.VariableModel, 0)

	handler := mapCreate.New(mapCreate.NewModel(projectId, name, entries), auth.NewTestingAuthentication(false, ""))

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID)

	gomega.Expect(name).Should(gomega.Equal(view.Name))

	return view
}

func testCreateList(projectId, name string) createList2.View {
	handler := createList2.New(createList2.NewModel(projectId, name, []createList2.Variable{}), auth.NewTestingAuthentication(false, ""))

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

	handler := addGroups.New(addGroups.NewModel(projectId, groups), auth.NewTestingAuthentication(false, projectId))
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

	model := addToMap.NewModel(projectId, name, variableModel, references, []string{})
	handler := addToMap.New(model, auth.NewTestingAuthentication(false, ""))

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

	model := addToList.NewModel(projectId, name, variableModel, references, []string{})
	handler := addToList.New(model, auth.NewTestingAuthentication(false, ""))

	view, err := handler.Handle()
	gomega.Expect(err).Should(gomega.BeNil())

	return view
}

func publishFullProject(projectId string) (addToMap.View, publish.View) {
	groups := testCreateGroups(projectId, 5)

	map1 := testCreateMap(projectId, "map1")

	referenceMap := testCreateMap(projectId, "referenceMap")
	referenceMapItem1 := testAddToMap(projectId, referenceMap.ID, "reference-map-1", []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
		return value.ID
	}))
	referenceMapItem2 := testAddToMap(projectId, referenceMap.ID, "reference-map-2", []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
		return value.ID
	}))

	referenceList := testCreateList(projectId, "referenceList")
	referenceListItem1 := testAddToList(projectId, referenceList.ID, "reference-list-1", []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
		return value.ID
	}))
	referenceListItem2 := testAddToList(projectId, referenceList.ID, "reference-list-2", []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
		return value.ID
	}))

	addToMapModel := testAddToMap(projectId, map1.ID, "mapItemName", []shared.Reference{
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
			StructureName: referenceList.Name,
			StructureType: "list",
			VariableID:    referenceListItem1.ID,
		},
		{
			Name:          "fourth",
			StructureName: referenceList.Name,
			StructureType: "list",
			VariableID:    referenceListItem2.ID,
		},
	}, sdk.Map(groups, func(idx int, value addGroups.View) string {
		return value.ID
	}))

	handler := publish.New(publish.NewModel(projectId, "v1"), auth.NewTestingAuthentication(false, ""))
	model, err := handler.Handle()
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(model.ID).ShouldNot(gomega.BeEmpty())
	gomega.Expect(model.Name).ShouldNot(gomega.BeEmpty())

	return addToMapModel, model
}
