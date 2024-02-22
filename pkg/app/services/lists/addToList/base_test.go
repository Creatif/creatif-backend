package addToList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain"
	"creatif/pkg/app/services/groups/addGroups"
	createList2 "creatif/pkg/app/services/lists/createList"
	"creatif/pkg/app/services/locales"
	createProject2 "creatif/pkg/app/services/projects/createProject"
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
	GinkgoRunSpecs(t, "Declaration -> CRUD tests")
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
	runLogger()

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

func testCreateList(projectId, name string) createList2.View {
	handler := createList2.New(createList2.NewModel(projectId, name, []createList2.Variable{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

	list, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(list.ID)

	gomega.Expect(list.Name).Should(gomega.Equal(name))

	return list
}

func testCreateGroups(projectId string, numOfGroups int) []string {
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
	gomega.Expect(err).Should(gomega.BeNil())

	return sdk.Map(model, func(idx int, value addGroups.View) string {
		return value.ID
	})
}

func testAddToList(projectId, name, variableName string, references []shared.Reference, groups []string) View {
	variableModel := VariableModel{
		Name:      variableName,
		Metadata:  nil,
		Groups:    groups,
		Value:     nil,
		Locale:    "eng",
		Behaviour: "modifiable",
	}

	model := NewModel(projectId, name, variableModel, references)
	handler := New(model, auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

	view, err := handler.Handle()
	gomega.Expect(err).Should(gomega.BeNil())

	return view
}
