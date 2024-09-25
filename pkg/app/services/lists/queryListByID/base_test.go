package queryListByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/lists/addToList"
	createList2 "creatif/pkg/app/services/lists/createList"
	"creatif/pkg/app/services/locales"
	createProject2 "creatif/pkg/app/services/projects/createProject"
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
	GinkgoRunSpecs(t, "Declaration -> CRUD tests")
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
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", domain.PUBLISHED_REFERENCES_TABLE))
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

func testCreateListAndReturnIds(projectId, name string, varNum int) (string, []map[string]string) {
	variables := make([]createList2.Variable, varNum)
	for i := 0; i < varNum; i++ {
		variables[i] = createList2.Variable{
			Name:      fmt.Sprintf("one-%d", i),
			Metadata:  nil,
			Groups:    nil,
			Locale:    "eng",
			Behaviour: "readonly",
			Value:     nil,
		}
	}

	handler := createList2.New(createList2.NewModel(projectId, name, variables), auth.NewTestingAuthentication(false, ""))

	list, err := handler.Handle()
	gomega.Expect(err).Should(gomega.BeNil())
	testAssertIDValid(list.ID)

	gomega.Expect(list.Name).Should(gomega.Equal(name))

	var savedVariables []declarations.ListVariable
	res := storage2.Gorm().Where("list_id = ?", list.ID).Find(&savedVariables)
	gomega.Expect(res.Error).Should(gomega.BeNil())

	return list.ID, sdk.Map(savedVariables, func(idx int, value declarations.ListVariable) map[string]string {
		return map[string]string{
			"id":   value.ID,
			"name": value.Name,
		}
	})
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
	gomega.Expect(err).Should(gomega.BeNil())

	return model
}

func testAddToList(projectId, name string, references []shared.Reference, groups []string) addToList.View {
	variableModel := addToList.VariableModel{
		Name:      fmt.Sprintf("new add variable"),
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
