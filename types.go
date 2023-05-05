package asyncify

type promiseState int

type promise struct {
	state     promiseState
	result    interface{}
	err       error
	thenFn    func(interface{}) interface{}
	catchFn   func(error) interface{}
	finallyFn func()
	awaitChan chan struct{}
}
