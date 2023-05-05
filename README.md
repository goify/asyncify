# Asyncify

![build](https://github.com/iamando/asyncify/workflows/build/badge.svg)
![license](https://img.shields.io/github/license/iamando/asyncify?color=success)
![Go version](https://img.shields.io/github/go-mod/go-version/iamando/asyncify)
[![GoDoc](https://godoc.org/github.com/iamando/asyncify?status.svg)](https://godoc.org/github.com/iamando/asyncify)

This is a Go module that provides an implementation of Promises, similar to those in JavaScript, including support for then, catch, and finally. It also allows for Await to be used to block the execution of the program until the promise resolves or rejects.

## Installation

To use this module in your Go project, run the following command:

```bash
go get github.com/iamando/asyncify
```

## Implementation

This Go module provides a simple implementation of Promises in JavaScript with support for `then`, `catch`, and `finally` handlers.

### Usage with `then`, `catch`, and `finally`

Here's an example of how to use the promise type with `then`, `catch`, and `finally`:

```go
func main() {
  promise := Promise(func(resolve func(interface{}), reject func(error)) {
    // Do some asynchronous operation here
    time.Sleep(2 * time.Second)

    if true {
      resolve("Success!")
    } else {
      reject(errors.New("Error!"))
    }
  })

  promise
    .Then(func(result interface{}) interface{} {
      // Handle fulfilled promise
      fmt.Println(result)
      return "Done!"
    })
    .Catch(func(err error) interface{} {
      // Handle rejected promise
      fmt.Println(err.Error())
      return "Done!"
    })
    .Finally(func() {
      // Handle either case
      fmt.Println("Finished!")
    })
}
```

In this example, a new `promise` is created with an asynchronous operation that takes 2 seconds to complete. The `then` handler is used to handle the successful resolution of the promise, the `catch` handler is used to handle any errors that may occur, and the `finally` handler is used to handle either case.

### Usage with `Await`

Here's an example of how to use the `Await` method to block the execution of the program until the promise resolves or rejects:

```go
func main() {
  promise := Promise(func(resolve func(interface{}), reject func(error)) {
    // Do some asynchronous operation here
    time.Sleep(2 * time.Second)

    if true {
      resolve("Success!")
    } else {
      reject(errors.New("Error!"))
    }
  })

  result, err := promise.Await()

  if err != nil {
    // Handle error
    fmt.Println(err.Error())
  } else {
    // Handle result
    fmt.Println(result)
  }
}
```

In this example, a new `promise` is created with an asynchronous operation that takes 2 seconds to complete. The `Await` method is used to block the execution of the program until the promise resolves or rejects. If the promise is resolved, the result is returned. If the promise is rejected, an error is returned.

## API

### type `PromiseStruct`

The `PromiseStruct` type represents a promise that will be resolved with a value or rejected with an error. It has the following methods:

### `Promise(executor func(resolve func(interface{}), reject func(error))) *PromiseStruct`

`Promise` creates a new promise with an executor function that takes two functions as arguments: `resolve` and `reject`. `resolve` should be called with the result of the promise when it is successfully resolved, and `reject` should be called with an error if the promise is rejected.

### `Then(fn func(interface{}) interface{}) *PromiseStruct`

`Then` creates a new promise that is resolved with the result of the `fn` function when the original promise is fulfilled. If the original promise is rejected, the new promise is rejected with the same error.

### `Catch(fn func(error) interface{}) *PromiseStruct`

`Catch` creates a new promise that is resolved with the result of the `fn` function when the original promise is rejected. If the original promise is fulfilled, the new promise is resolved with the same result.

### `Finally(fn func()) *PromiseStruct`

`Finally` creates a new promise

### `Await() (interface{}, error)`

The `Await` function is a blocking function that allows you to wait for a promise to resolve or reject. It waits until the promise state changes from `pending` to either `fulfilled` or `rejected`. If the promise is fulfilled, `Await` returns nil, and if the promise is rejected, it returns the error that caused the rejection.

## Testing

```bash
go test
```

## Support

Asyncify is an MIT-licensed open source project. It can grow thanks to the sponsors and support.

## License

Asyncify is [MIT licensed](LICENSE).
