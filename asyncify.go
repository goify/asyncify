package asyncify

// Resolve returns a new Promise that is already resolved with the given value.
// The returned Promise's `result` channel will receive the value immediately
// after the function call. The `err` channel will not be used.
func Resolve(val interface{}) *Promise {
	p := &Promise{
		result: make(chan interface{}),
		err:    make(chan error),
	}

	go func() {
		p.result <- val
	}()

	return p
}
