package workloadEngine

type Engine interface {
	Start()
}

type engine[T any] struct {
	jobs []T
}

func (e engine[T]) Start() {
	l := len(e.jobs) / 20
	r := len(e.jobs) % 20
}

func NewEngine[T any](jobs []T) {

}
