package asyncify

type Promise struct {
	result chan interface{}
	err    chan error
}
