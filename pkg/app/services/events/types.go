package events

type DispachableEvent interface {
	Type() string
	Data() []byte
}
