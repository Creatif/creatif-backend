package getProjectMetadata

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain"
	createList2 "creatif/pkg/app/services/lists/createList"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/maps/mapCreate"
	"creatif/pkg/app/services/projects/createProject"
	createVariable2 "creatif/pkg/app/services/variables/createVariable"
	"creatif/pkg/lib/logger"
	storage2 "creatif/pkg/lib/storage"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"log"
	"os"
	"testing"
	"time"
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
	GinkgoRunSpecs(t, "Project -> CRUD tests")
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

	gomega.Expect(locales.Store()).Should(gomega.BeNil())

	if err != nil {
		log.Fatalln(err)
	}
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
	res := storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.VARIABLES_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", domain.PROJECT_TABLE))
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
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", domain.GROUPS_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", domain.VARIABLE_GROUPS_TABLE))
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

func testCreateProject(name string) auth.Authentication {
	a := auth.NewTestingAuthentication(true, "")

	handler := createProject.New(createProject.NewModel(name), a, logger.NewLogBuilder())

	model, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(model.ID)

	f := auth.NewFrontendTestingAuthentication(auth.AuthenticatedUser{
		ID:        a.User().ID,
		Name:      a.User().Name,
		LastName:  a.User().LastName,
		Email:     a.User().Email,
		Refresh:   a.User().Refresh,
		ProjectID: model.ID,
		ApiKey:    "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	})

	gomega.Expect(model.Name).Should(gomega.Equal(name))

	return f
}

func testCreateList(authentication auth.Authentication, name string, varNum int) createList2.View {
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

	handler := createList2.New(createList2.NewModel(authentication.User().ProjectID, name, variables), authentication, logger.NewLogBuilder())

	list, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(list.ID)

	gomega.Expect(list.Name).Should(gomega.Equal(name))

	return list
}

func testCreateMap(authentication auth.Authentication, name string, variablesNum int) mapCreate.View {
	entries := make([]mapCreate.VariableModel, 0)
	fragmentedGroups := map[string]int{}
	fragmentedGroups["one"] = 0
	fragmentedGroups["two"] = 0
	fragmentedGroups["three"] = 0

	m := map[string]interface{}{
		"one":   "one",
		"two":   []string{"one", "two", "three"},
		"three": []int{1, 2, 3},
		"four":  453,
	}

	b, err := json.Marshal(m)
	gomega.Expect(err).Should(gomega.BeNil())

	for i := 0; i < variablesNum; i++ {
		var value interface{}
		value = "my value"
		if i%2 == 0 {
			value = true
		}

		if i%3 == 0 {
			value = map[string]interface{}{
				"one":   "one",
				"two":   []string{"one", "two", "three"},
				"three": []int{1, 2, 3},
				"four":  453,
			}
		}

		var groups []string = []string{"unfragmented"}
		if i%2 == 0 {
			groups = append(groups, "one")
			fragmentedGroups["one"]++
		}

		if i%3 == 0 {
			groups = append(groups, "two")
			fragmentedGroups["two"]++
		}

		if i%5 == 0 {
			groups = append(groups, "three")
			fragmentedGroups["three"]++
		}

		v, err := json.Marshal(value)
		gomega.Expect(err).Should(gomega.BeNil())

		variableModel := mapCreate.VariableModel{
			Name:      fmt.Sprintf("name-%d", i),
			Metadata:  b,
			Groups:    groups,
			Value:     v,
			Behaviour: "modifiable",
			Locale:    "eng",
		}

		entries = append(entries, variableModel)
	}

	handler := mapCreate.New(mapCreate.NewModel(authentication.User().ProjectID, name, entries), authentication, logger.NewLogBuilder())

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID)

	gomega.Expect(name).Should(gomega.Equal(view.Name))
	gomega.Expect(len(view.Variables)).Should(gomega.Equal(variablesNum))

	return view
}

func testCreateDetailedVariable(authentication auth.Authentication, locale, name, behaviour string, groups []string, metadata []byte) createVariable2.View {
	b, _ := json.Marshal(map[string]interface{}{
		"one":  1,
		"two":  "three",
		"four": "six",
	})

	handler := createVariable2.New(createVariable2.NewModel(authentication.User().ProjectID, locale, name, behaviour, groups, metadata, b), authentication, logger.NewLogBuilder())

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID)

	return view
}
