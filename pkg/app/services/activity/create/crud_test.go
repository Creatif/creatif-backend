package create

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/datatypes"
)

var _ = ginkgo.Describe("Activity", func() {
	ginkgo.It("should be created with arbitrary data", ginkgo.Label("activity"), func() {
		ginkgo.Skip("")
		projectId := testCreateProject("project")

		data := map[string]string{
			"type":  "visit",
			"path":  "path/to/link",
			"title": "You visited groups",
		}

		b, err := json.Marshal(data)
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, b), auth.NewTestingAuthentication(false, ""))

		model, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(model.ID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.CreatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(model.Data).ShouldNot(gomega.BeNil())

		savedData := make(map[string]string)
		gomega.Expect(json.Unmarshal(model.Data.(datatypes.JSON), &savedData)).Should(gomega.BeNil())

		gomega.Expect(savedData["type"]).Should(gomega.Equal("visit"))
		gomega.Expect(savedData["path"]).Should(gomega.Equal("path/to/link"))
		gomega.Expect(savedData["title"]).Should(gomega.Equal("You visited groups"))
	})

	ginkgo.It("should be called successively more than 10 times but only one should be written since it is a visit", ginkgo.Label("activity"), func() {
		ginkgo.Skip("")

		projectId := testCreateProject("project")

		data := map[string]string{
			"type":  "visit",
			"path":  "path/to/link",
			"title": "You visited groups",
		}

		b, err := json.Marshal(data)
		gomega.Expect(err).Should(gomega.BeNil())

		for i := 0; i < 100; i++ {
			handler := New(NewModel(projectId, b), auth.NewTestingAuthentication(false, ""))

			_, err := handler.Handle()
			gomega.Expect(err).Should(gomega.BeNil())
		}

		sql := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE project_id = ?", (app.Activity{}).TableName())

		var count int
		res := storage.Gorm().Raw(sql, projectId).Scan(&count)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		gomega.Expect(count).Should(gomega.Equal(1))
	})
})
