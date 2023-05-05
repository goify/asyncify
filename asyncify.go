package asyncify

// Resolve returns a new Promise that is already resolved with the given value.
// The returned Promise's `result` channel will receive the value immediately
// after the function call. The `err` channel will not be used.
func Resolve(val interface{}) *Promise {
	promise := &Promise{
		result: make(chan interface{}),
		err:    make(chan error),
	}

	go func() {
		promise.result <- val
	}()

	return promise
}

// Reject returns a new Promise that is already rejected with the given error.
// The returned Promise's `err` channel will receive the error immediately
// after the function call. The `result` channel will not be used.
func Reject(err error) *Promise {
	promise := &Promise{
		result: make(chan interface{}),
		err:    make(chan error),
	}

	go func() {
		promise.err <- err
	}()

	return promise
}

// Then appends a new callback function to the promise's `callbacks` list.
// The callback function is invoked when the promise is resolved, and the
// value returned by the callback is used to resolve the new promise that
// is returned by this function. If the callback throws an error, the new
// promise is rejected with the error.
func (p *Promise) Then(fn func(interface{}) interface{}) *Promise {
	promise := &Promise{
		result: make(chan interface{}),
		err:    make(chan error),
	}

	go func() {
		select {
		case val := <-p.result:
			res := fn(val)

			if p, ok := res.(*Promise); ok {
				promise.result = p.result
				promise.err = p.err
			} else {
				promise.result <- res
			}

		case err := <-p.err:
			promise.err <- err
		}
	}()

	return promise
}
