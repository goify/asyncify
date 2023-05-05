package asyncify

// Promise creates and returns a new Promise object that represents the eventual completion
// (or failure) of an asynchronous operation, and executes the specified executor function
// immediately. The executor function takes two arguments, a resolve function and a reject function,
// that allow the promise to be resolved with a value or rejected with a reason, respectively.
// The Promise object returned by this function has methods `Then`, `Catch`, and `Finally`
// that allow for chaining and handling of the eventual fulfillment or rejection of the promise.
// The Promise object also provides a blocking `Await` method that can be used to wait for the
// promise to be resolved or rejected.
func Promise(executor func(resolve func(interface{}), reject func(error))) *PromiseStruct {
	promise := &PromiseStruct{
		state:     pending,
		awaitChan: make(chan struct{}),
	}

	resolve := func(result interface{}) {
		if promise.state != pending {
			return
		}

		if promise.thenFn != nil {
			promise.result = promise.thenFn(result)
		} else {
			promise.result = result
		}

		promise.state = fulfilled

		if promise.finallyFn != nil {
			promise.finallyFn()
		}

		close(promise.awaitChan)
	}

	reject := func(err error) {
		if promise.state != pending {
			return
		}

		if promise.catchFn != nil {
			promise.result = promise.catchFn(err)
		} else {
			promise.err = err
		}

		promise.state = rejected

		if promise.finallyFn != nil {
			promise.finallyFn()
		}

		close(promise.awaitChan)
	}

	go executor(resolve, reject)

	return promise
}

// Then registers a callback function to be called when the promise is resolved. If the promise is already resolved, the callback is called immediately with the resolved value. If the promise is rejected, the callback is skipped.
// The callback function takes one argument, the resolved value of the promise, and should return a value or a new promise that will be resolved with that value.
// Returns a new promise that is resolved with the return value of the callback function or rejected with the same reason as the original promise, if the callback function throws an error.
func (p *PromiseStruct) Then(fn func(interface{}) interface{}) *PromiseStruct {
	promise := &PromiseStruct{
		state:     pending,
		awaitChan: make(chan struct{}),
	}

	p.thenFn = fn

	if p.state == fulfilled {
		if p.thenFn != nil {
			promise.result = p.thenFn(p.result)
		} else {
			promise.result = p.result
		}

		promise.state = fulfilled

		if p.finallyFn != nil {
			p.finallyFn()
		}

		close(promise.awaitChan)

	} else if p.state == rejected {
		if p.catchFn != nil {
			promise.result = p.catchFn(p.err)
		} else {
			promise.err = p.err
		}

		promise.state = rejected

		if p.finallyFn != nil {
			p.finallyFn()
		}

		close(promise.awaitChan)
	}

	return promise
}

// Catch registers a callback function to be called when the promise is rejected. If the promise is already rejected, the callback is called immediately with the rejection reason. If the promise is resolved, the callback is skipped.
// The callback function takes one argument, the rejection reason of the promise, and should return a value or a new promise that will be resolved with that value.
// Returns a new promise that is resolved with the return value of the callback function or rejected with the same reason as the original promise, if the callback function throws an error.
func (p *PromiseStruct) Catch(fn func(error) interface{}) *PromiseStruct {
	promise := &PromiseStruct{
		state:     pending,
		awaitChan: make(chan struct{}),
	}

	p.catchFn = fn

	if p.state == fulfilled {
		promise.result = p.result
		promise.state = fulfilled

		if p.finallyFn != nil {
			p.finallyFn()
		}

		close(promise.awaitChan)
	} else if p.state == rejected {
		if p.catchFn != nil {
			promise.result = p.catchFn(p.err)
		} else {
			promise.err = p.err
		}

		promise.state = rejected

		if p.finallyFn != nil {
			p.finallyFn()
		}

		close(promise.awaitChan)
	}

	return promise
}

// Finally registers a callback function to be called when the promise is either resolved or rejected. If the promise is already resolved or rejected, the callback is called immediately.
// The callback function takes no arguments and should not return anything.
// Returns the same promise instance to allow for chaining of methods.
func (promise *PromiseStruct) Finally(fn func()) *PromiseStruct {
	promise.finallyFn = fn

	if promise.state != pending {
		promise.finallyFn()
	}

	return promise
}

// Await blocks the execution of the program until the promise resolves or rejects,
// and returns either the resolved value or an error.
// It returns an error only if the promise was rejected, and the resolved value otherwise.
func (promise *PromiseStruct) Await() (interface{}, error) {
	<-promise.awaitChan

	if promise.state == rejected {
		return nil, promise.err
	}

	return promise.result, nil
}
