package services

import "encoding/json"

type ValueSerializer[T any] struct {
	value string
}

func (vs ValueSerializer[T]) Serialize() []byte {
	return []byte(vs.value)
}

func (vs ValueSerializer[T]) UnSerialize(model interface{}) (T, error) {
	var val T
	return val, json.Unmarshal([]byte(vs.value), &model)
}
