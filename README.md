# backoff

`backoff` is a Go package that provides iterators powered by
[github.com/cenkalti/backoff/v4](https://pkg.go.dev/github.com/cenkalti/backoff/v4).

It offers a flexible and easy-to-use way to implement various backoff strategies
in your Go applications.

## Features

- Easy integration with Go's `for` range statements
- Supports multiple backoff policies:
  - Stop BackOff
  - Zero BackOff
  - Constant BackOff
  - Exponential BackOff
- Context-aware for graceful cancellation
- Customizable options for fine-tuning backoff behavior
- All of BackOff's features come from [github.com/cenkalti/backoff/v4](https://pkg.go.dev/github.com/cenkalti/backoff/v4), which is used by many.

## Installation

To install `backoff`, use `go get`:

```
go get github.com/goaux/backoff
```

## Usage

Here's a basic example of how to use `backoff`:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/goaux/backoff"
)

func main() {
    exponentialBackOff := backoff.NewExponential(
        backoff.WithInitialInterval(100 * time.Millisecond),
        backoff.WithMaxRetries(5),
    )
    ctx := context.Background()
    for i := range exponentialBackOff(ctx) {
        fmt.Printf("Attempt %d\n", i)
        // Perform your operation here
        // Break the loop if the operation succeeds
    }
}
```

## License

This project is licensed under the [Apache-2.0](LICENSE).

## Special Thanks

@cenkalti https://github.com/cenkalti/backoff
