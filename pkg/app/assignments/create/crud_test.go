package create

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Assignment CRUD success test", func() {
	ginkgo.It("should create an assignment text node when the node does not exists", ginkgo.Label("assignment", "crud", "success", "1"), func() {
		name := uuid.NewString()
		declarationNode := testCreateBasicDeclarationTextNode(name, "modifiable")

		b, _ := json.Marshal("this is a text node")

		handler := New(NewCreateNodeModel(declarationNode.Name, b))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)
	})

	ginkgo.It("should update an assignment text node when the node already exists", ginkgo.Label("assignment", "crud", "success", "2"), func() {
		name := uuid.NewString()
		testCreateBasicAssignmentTextNode(name)
		text := "this is a changed text value"

		b, _ := json.Marshal(text)

		handler := New(NewCreateNodeModel(name, b))

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())
		var textNode assignments.NodeText
		gomega.Expect(storage.Gorm().Where("assignment_node_id = ?", assignmentNode.ID).First(&textNode).Error).Should(gomega.BeNil())

		var changedText string
		gomega.Expect(json.Unmarshal(textNode.Value, &changedText))
		gomega.Expect(changedText).Should(gomega.Equal(text))
	})
})
