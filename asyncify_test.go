package asyncify

import (
	"errors"
	"testing"
	"time"
)

func TestPromise(t *testing.T) {
	// Create a new promise
	promise := Promise(func(resolve func(interface{}), reject func(error)) {
		// Simulate an asynchronous operation
		go func() {
			// Resolve the promise with a value
			resolve("hello")
		}()
	})

	// Test the `Then` method
	promise.Then(func(result interface{}) interface{} {
		// Test that the result is correct
		if result != "hello" {
			t.Errorf("Expected result to be 'hello', but got '%v'", result)
		}

		// Return a new value
		return "world"
	}).Then(func(result interface{}) interface{} {
		// Test that the result is correct
		if result != "world" {
			t.Errorf("Expected result to be 'world', but got '%v'", result)
		}

		// Return an error
		return errors.New("oops")
	}).Catch(func(err error) interface{} {
		// Test that the error is correct
		if err.Error() != "oops" {
			t.Errorf("Expected error to be 'oops', but got '%v'", err.Error())
		}

		// Return a new value
		return "caught"
	}).Finally(func() {
		// Test that the promise is fulfilled
		if promise.state != fulfilled {
			t.Error("Expected promise to be fulfilled")
		}

		// Close the test channel
		close(promise.awaitChan)
	})

	// Wait for the promise to resolve
	<-promise.awaitChan
}

func TestPromiseAwait(t *testing.T) {
	// Create a new promise that resolves after 100 milliseconds
	p := Promise(func(resolve func(interface{}), reject func(error)) {
		time.Sleep(100 * time.Millisecond)
		resolve("Hello, world!")
	})

	// Call Await on the promise and capture the result
	result, err := p.Await()

	// Check that there is no error
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check that the result matches the expected value
	expected := "Hello, world!"
	if result != expected {
		t.Errorf("Expected %v, but got: %v", expected, result)
	}
}
