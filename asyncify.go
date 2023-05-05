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

// Reject returns a new Promise that is already rejected with the given error.
// The returned Promise's `err` channel will receive the error immediately
// after the function call. The `result` channel will not be used.
func Reject(err error) *Promise {
	p := &Promise{
		result: make(chan interface{}),
		err:    make(chan error),
	}

	go func() {
		p.err <- err
	}()

	return p
}
