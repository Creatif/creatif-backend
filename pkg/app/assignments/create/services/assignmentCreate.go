package services

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssignmentCreate struct {
	nodeName          string
	value             []byte
	declarationNodeID uuid.UUID
}

type AssignmentCreateResult struct {
	Node  assignments.Node
	Value interface{}
}

func NewAssignmentCreate(nodeName string, value []byte, declarationNodeID uuid.UUID) AssignmentCreate {
	return AssignmentCreate{
		value:             value,
		declarationNodeID: declarationNodeID,
		nodeName:          nodeName,
	}
}

func (a AssignmentCreate) CreateOrUpdate() (AssignmentCreateResult, error) {
	model := assignments.NewNode(a.nodeName, a.declarationNodeID)

	var exists assignments.Node
	res := storage.Gorm().Where("name = ?", a.nodeName).First(&exists)

	// any other error than record not exists if a failure and processing should not continue
	if exists.ID.String() == "" && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return AssignmentCreateResult{}, res.Error
	}

	var createdOrUpdatedValue []byte

	if err := storage.Transaction(func(tx *gorm.DB) error {
		// record exists, update the value node
		if exists.ID.String() != "" {
			if res := tx.Table((&assignments.ValueNode{}).TableName()).Where("assignment_node_id", exists.ID).Update("value", a.value); res.Error != nil {
				return res.Error
			}

			createdOrUpdatedValue = a.value
		}

		// record does not exist, create assignment node and value node
		if exists.ID.String() == "" {
			node := assignments.NewNode(a.nodeName, a.declarationNodeID)
			if res := tx.Create(&node); res != nil {
				return res.Error
			}

			valueNode := assignments.NewValueNode(node.ID, a.value)
			if res := tx.Create(&valueNode); res.Error != nil {
				return res.Error
			}

			createdOrUpdatedValue = valueNode.Value
		}

		return nil
	}); err != nil {
		return a.error(err)
	}

	var v interface{}
	if createdOrUpdatedValue != nil {
		if err := json.Unmarshal(createdOrUpdatedValue, &v); err != nil {
			return AssignmentCreateResult{}, appErrors.NewDatabaseError(err).AddError("Node.Create.Service.AssignmentCreate", nil)
		}
	} else {
		v = createdOrUpdatedValue
	}

	if exists.ID.String() != "" {
		return AssignmentCreateResult{
			Node:  exists,
			Value: v,
		}, nil
	}

	return AssignmentCreateResult{
		Node:  model,
		Value: v,
	}, nil
}

func (a AssignmentCreate) error(err error) (AssignmentCreateResult, error) {
	return AssignmentCreateResult{}, appErrors.NewDatabaseError(err).AddError("AssignmentCreate", nil)
}
