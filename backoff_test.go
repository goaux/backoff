package backoff_test

import (
	"context"
	"fmt"
	"time"

	"github.com/goaux/backoff"
)

func ExampleNewConstant_stopBackOff() {
	ctx := context.TODO()
	unit := 100 * time.Millisecond

	// interval < 0
	// [StopBackOff](https://pkg.go.dev/github.com/cenkalti/backoff/v4#StopBackOff)
	constantBackOff := backoff.NewConstant(-1, backoff.WithMaxRetries(3))
	start := time.Now()
	for i := range constantBackOff(ctx) {
		fmt.Println("stop#false", i, (time.Since(start) / unit * unit).String())
		if false {
			break
		}
	}
	start = time.Now()
	for i := range constantBackOff(ctx) {
		fmt.Println("stop#true", i, (time.Since(start) / unit * unit).String())
		if true {
			break
		}
	}
	// Output:
	// stop#false 0 0s
	// stop#true 0 0s
}

func ExampleNewConstant_zeroBackOff() {
	ctx := context.TODO()
	unit := 100 * time.Millisecond

	// interval == 0
	// [ZeroBackOff](https://pkg.go.dev/github.com/cenkalti/backoff/v4#ZeroBackOff)
	constantBackOff := backoff.NewConstant(0, backoff.WithMaxRetries(4))
	start := time.Now()
	for i := range constantBackOff(ctx) {
		fmt.Println("zero#false", i, (time.Since(start) / unit * unit).String())
		if false {
			break
		}
	}
	start = time.Now()
	for i := range constantBackOff(ctx) {
		fmt.Println("zero#true", i, (time.Since(start) / unit * unit).String())
		if true {
			break
		}
	}
	// Output:
	// zero#false 0 0s
	// zero#false 1 0s
	// zero#false 2 0s
	// zero#false 3 0s
	// zero#false 4 0s
	// zero#true 0 0s
}

func ExampleNewConstant() {
	ctx := context.TODO()
	unit := 100 * time.Millisecond

	// interval > 0
	// [ConstantBackOff](https://pkg.go.dev/github.com/cenkalti/backoff/v4#ConstantBackOff)
	constant := backoff.NewConstant(2*unit, backoff.WithMaxRetries(5))
	start := time.Now()
	for i := range constant(ctx) {
		fmt.Println("constant#false", i, (time.Since(start) / unit * unit).String())
		if false {
			break
		}
	}
	start = time.Now()
	for i := range constant(ctx) {
		fmt.Println("constant#true", i, (time.Since(start) / unit * unit).String())
		if true {
			break
		}
	}
	// Output:
	// constant#false 0 0s
	// constant#false 1 200ms
	// constant#false 2 400ms
	// constant#false 3 600ms
	// constant#false 4 800ms
	// constant#false 5 1s
	// constant#true 0 0s
}

func ExampleNewConstant_cancelContext() {
	ctx := context.TODO()
	unit := 100 * time.Millisecond

	// interval > 0
	// [ConstantBackOff](https://pkg.go.dev/github.com/cenkalti/backoff/v4#ConstantBackOff)
	constantBackOff := backoff.NewConstant(2*unit, backoff.WithMaxRetries(5))
	start := time.Now()
	ctx, cancel := context.WithTimeout(ctx, 3*unit)
	defer cancel()
	for i := range constantBackOff(ctx) {
		fmt.Println("constant#false", i, (time.Since(start) / unit * unit).String())
		if false {
			break
		}
	}
	fmt.Println("cancel", (time.Since(start) / unit * unit).String())
	// Output:
	// constant#false 0 0s
	// constant#false 1 200ms
	// cancel 300ms
}

func ExampleNewExponential() {
	ctx := context.TODO()
	unit := 50 * time.Millisecond

	// [ExponentialBackOff](https://pkg.go.dev/github.com/cenkalti/backoff/v4#ExponentialBackOff)
	exponentialBackOff := backoff.NewExponential(
		backoff.WithInitialInterval(2*unit),
		backoff.WithRandomizationFactor(0), // Set to 0 for demonstrative purpose.
		backoff.WithMaxRetries(7),
	)
	start := time.Now()
	for i := range exponentialBackOff(ctx) {
		fmt.Println("exponential#false", i, (time.Since(start) / unit * unit).String())
		if false {
			break
		}
	}
	start = time.Now()
	for i := range exponentialBackOff(ctx) {
		fmt.Println("exponential#true", i, (time.Since(start) / unit * unit).String())
		if true {
			break
		}
	}
	// Output:
	// exponential#false 0 0s
	// exponential#false 1 100ms
	// exponential#false 2 250ms
	// exponential#false 3 450ms
	// exponential#false 4 800ms
	// exponential#false 5 1.3s
	// exponential#false 6 2.05s
	// exponential#false 7 3.2s
	// exponential#true 0 0s
}

func ExampleNewExponential_cancelContext() {
	ctx := context.TODO()
	unit := 50 * time.Millisecond

	// [ExponentialBackOff](https://pkg.go.dev/github.com/cenkalti/backoff/v4#ExponentialBackOff)
	exponentialBackOff := backoff.NewExponential(
		backoff.WithInitialInterval(2*unit),
		backoff.WithRandomizationFactor(0), // Set to 0 for demonstrative purpose.
		backoff.WithMaxRetries(7),
	)
	start := time.Now()
	ctx, cancel := context.WithTimeout(ctx, 3*unit)
	defer cancel()
	for i := range exponentialBackOff(ctx) {
		fmt.Println("exponential#false", i, (time.Since(start) / unit * unit).String())
		if false {
			break
		}
	}
	fmt.Println("cancel", (time.Since(start) / unit * unit).String())
	// Output:
	// exponential#false 0 0s
	// exponential#false 1 100ms
	// cancel 150ms
}

func Example() {
	ctx := context.TODO()
	unit := 100 * time.Millisecond

	// [ExponentialBackOff](https://pkg.go.dev/github.com/cenkalti/backoff/v4#ExponentialBackOff)
	exponentialBackOff := backoff.NewExponential(
		backoff.WithInitialInterval(unit),
		backoff.WithRandomizationFactor(0), // Set to 0 for demonstrative purpose.
		backoff.WithMaxRetries(7),
	)
	start := time.Now()
	for i := range exponentialBackOff(ctx) {
		fmt.Println("task", i, (time.Since(start) / unit * unit).String())
		if false { // Here we assume the task will always fail, so no break will occur.
			break
		}
	}
	// Output:
	// task 0 0s
	// task 1 100ms
	// task 2 200ms
	// task 3 400ms
	// task 4 800ms
	// task 5 1.3s
	// task 6 2s
	// task 7 3.2s
}
