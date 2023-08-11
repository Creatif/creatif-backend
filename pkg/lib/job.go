package pkg

type Job[T any, F any, K any] interface {
	Validate() error
	Authenticate() error
	Authorize() error
	Logic() (K, error)
	Handle() (F, error)
}
