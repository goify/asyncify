package asyncify

// Promise creates and returns a new Promise object that represents the eventual completion
// (or failure) of an asynchronous operation, and executes the specified executor function
// immediately. The executor function takes two arguments, a resolve function and a reject function,
// that allow the promise to be resolved with a value or rejected with a reason, respectively.
// The Promise object returned by this function has methods `Then`, `Catch`, and `Finally`
// that allow for chaining and handling of the eventual fulfillment or rejection of the promise.
// The Promise object also provides a blocking `Await` method that can be used to wait for the
// promise to be resolved or rejected.
func Promise(executor func(resolve func(interface{}), reject func(error))) *promise {
	p := &promise{
		state:     pending,
		awaitChan: make(chan struct{}),
	}

	resolve := func(result interface{}) {
		if p.state != pending {
			return
		}

		if p.thenFn != nil {
			p.result = p.thenFn(result)
		} else {
			p.result = result
		}

		p.state = fulfilled

		if p.finallyFn != nil {
			p.finallyFn()
		}

		close(p.awaitChan)
	}

	reject := func(err error) {
		if p.state != pending {
			return
		}

		if p.catchFn != nil {
			p.result = p.catchFn(err)
		} else {
			p.err = err
		}

		p.state = rejected

		if p.finallyFn != nil {
			p.finallyFn()
		}

		close(p.awaitChan)
	}

	go executor(resolve, reject)

	return p
}

// Then registers a callback function to be called when the promise is resolved. If the promise is already resolved, the callback is called immediately with the resolved value. If the promise is rejected, the callback is skipped.
// The callback function takes one argument, the resolved value of the promise, and should return a value or a new promise that will be resolved with that value.
// Returns a new promise that is resolved with the return value of the callback function or rejected with the same reason as the original promise, if the callback function throws an error.
func (p *promise) Then(fn func(interface{}) interface{}) *promise {
	Promise := &promise{
		state:     pending,
		awaitChan: make(chan struct{}),
	}

	p.thenFn = fn

	if p.state == fulfilled {
		if p.thenFn != nil {
			Promise.result = p.thenFn(p.result)
		} else {
			Promise.result = p.result
		}

		Promise.state = fulfilled

		if p.finallyFn != nil {
			p.finallyFn()
		}

		close(Promise.awaitChan)
	} else if p.state == rejected {
		if p.catchFn != nil {
			Promise.result = p.catchFn(p.err)
		} else {
			Promise.err = p.err
		}

		Promise.state = rejected

		if p.finallyFn != nil {
			p.finallyFn()
		}

		close(Promise.awaitChan)
	}

	return Promise
}

// Catch registers a callback function to be called when the promise is rejected. If the promise is already rejected, the callback is called immediately with the rejection reason. If the promise is resolved, the callback is skipped.
// The callback function takes one argument, the rejection reason of the promise, and should return a value or a new promise that will be resolved with that value.
// Returns a new promise that is resolved with the return value of the callback function or rejected with the same reason as the original promise, if the callback function throws an error.
func (p *promise) Catch(fn func(error) interface{}) *promise {
	Promise := &promise{
		state:     pending,
		awaitChan: make(chan struct{}),
	}

	p.catchFn = fn

	if p.state == fulfilled {
		Promise.result = p.result
		Promise.state = fulfilled

		if p.finallyFn != nil {
			p.finallyFn()
		}

		close(Promise.awaitChan)
	} else if p.state == rejected {
		if p.catchFn != nil {
			Promise.result = p.catchFn(p.err)
		} else {
			Promise.err = p.err
		}

		Promise.state = rejected

		if p.finallyFn != nil {
			p.finallyFn()
		}

		close(Promise.awaitChan)
	}

	return Promise
}

// Finally registers a callback function to be called when the promise is either resolved or rejected. If the promise is already resolved or rejected, the callback is called immediately.
// The callback function takes no arguments and should not return anything.
// Returns the same promise instance to allow for chaining of methods.
func (p *promise) Finally(fn func()) *promise {
	p.finallyFn = fn

	if p.state != pending {
		p.finallyFn()
	}

	return p
}
