package services

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"gorm.io/gorm"
)

type AssignmentCreate struct {
	nodeName          string
	incomingValueType string
	value             interface{}
	declarationNodeID string
}

type AssignmentCreateResult struct {
	Node      assignments.Node
	Value     interface{}
	ValueType string
}

func NewAssignmentCreate(nodeName, valueType string, value interface{}, declarationNodeID string) AssignmentCreate {
	return AssignmentCreate{
		value:             value,
		declarationNodeID: declarationNodeID,
		nodeName:          nodeName,
		incomingValueType: valueType,
	}
}

/*
*
 1. Assignment node not exists, create new node together with the value node
 2. Assignment node exists:
    2.1. Check if the incoming value type is different from the current value type. If so, delete the value type
    2.2 If incoming value type and existing value type are the same, only update the value type
    2.3 If the incoming vlaue type and existing value type are not the same, create a new value type (the old one was already
    deleted in 2.1)
*/
func (a AssignmentCreate) CreateOrUpdate() (AssignmentCreateResult, error) {
	model := assignments.NewNode(a.nodeName, a.declarationNodeID)
	nodeTextValue, tOk := a.value.([]byte)
	nodeBooleanValue, bOk := a.value.(bool)
	var nextValue interface{}
	nextValueType := ""
	if tOk {
		nextValueType = assignments.ValueTextType
		nextValue = nodeTextValue
	} else if bOk {
		nextValueType = assignments.ValueBooleanType
		nextValue = nodeBooleanValue
	}

	var exists assignments.Node
	res := storage.Gorm().Where("name = ?", a.nodeName).First(&exists)

	// any other error than record not exists if a failure and processing should not continue
	if exists.ID == "" && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return AssignmentCreateResult{}, res.Error
	}

	if err := storage.Transaction(func(tx *gorm.DB) error {
		// record exists
		if exists.ID != "" {
			// if the incoming type and current value type are different, remove the difference
			if exists.ValueType == assignments.ValueTextType && bOk {
				if err := storage.DeleteBy((assignments.NodeText{}).TableName(), "assignment_node_id", exists.ID, &assignments.NodeText{}); err != nil {
					return err
				}
			}

			// if the incoming type and current value type are different, remove the difference
			if exists.ValueType == assignments.ValueBooleanType && tOk {
				if err := storage.DeleteBy((assignments.NodeBoolean{}).TableName(), "assignment_node_id", exists.ID, &assignments.NodeBoolean{}); err != nil {
					return err
				}
			}

			// if the incoming type and current value type are the same, only update the value type
			if a.incomingValueType == exists.ValueType {
				if tOk {
					var node assignments.NodeText
					if err := storage.GetBy(node.TableName(), "assignment_node_id", exists.ID, &node); err != nil {
						return err
					}

					node.Value = nodeTextValue
					if err := storage.Update(node.TableName(), &node); err != nil {
						return err
					}
				}

				if bOk {
					var node assignments.NodeBoolean
					if err := storage.GetBy(node.TableName(), "assignment_node_id", exists.ID, &node); err != nil {
						return err
					}

					node.Value = nodeBooleanValue
					if err := storage.Update(node.TableName(), &node); err != nil {
						return err
					}
				}
				// if the are different, create a new value node, existing value node is already deleted
			} else {
				exists.ValueType = nextValueType

				if tOk {
					valueNode := assignments.NewNodeText(exists.ID, nodeTextValue)
					if err := storage.Create(valueNode.TableName(), &valueNode, false); err != nil {
						return err
					}
				}

				if bOk {
					valueNode := assignments.NewNodeBoolean(exists.ID, nodeBooleanValue)
					if err := storage.Create(valueNode.TableName(), &valueNode, false); err != nil {
						return err
					}
				}

				if err := storage.Update(exists.TableName(), &exists); err != nil {
					return err
				}
			}

			return nil
		}

		// record does not exist
		if exists.ID == "" {
			model.ValueType = nextValueType

			if err := storage.Create(model.TableName(), &model, false); err != nil {
				return err
			}

			if tOk {
				node := assignments.NewNodeText(model.ID, nodeTextValue)
				if err := storage.Create(node.TableName(), &node, false); err != nil {
					return err
				}
			} else if bOk {
				node := assignments.NewNodeBoolean(model.ID, nodeBooleanValue)
				if err := storage.Create(node.TableName(), &node, false); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return a.error(err)
	}

	if exists.ID != "" {
		return AssignmentCreateResult{
			Node:      exists,
			Value:     nextValue,
			ValueType: nextValueType,
		}, nil
	}

	return AssignmentCreateResult{
		Node:      model,
		Value:     nextValue,
		ValueType: nextValueType,
	}, nil
}

func (a AssignmentCreate) error(err error) (AssignmentCreateResult, error) {
	return AssignmentCreateResult{}, appErrors.NewDatabaseError(err).AddError("AssignmentCreate", nil)
}
