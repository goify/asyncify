package asyncify

type PromiseState int

type PromiseStruct struct {
	state     PromiseState
	result    interface{}
	err       error
	thenFn    func(interface{}) interface{}
	catchFn   func(error) interface{}
	finallyFn func()
	awaitChan chan struct{}
}
