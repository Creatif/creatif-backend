package events

type DispachableEvent interface {
	Type() string
	Project() string
	Data() []byte
}
