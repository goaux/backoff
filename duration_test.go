package backoff_test

import (
	"fmt"
	"iter"
	"slices"
	"testing"
	"time"

	"github.com/goaux/backoff"
)

func ExampleNewConstantDuration() {
	backOffDuration := backoff.NewConstantDuration(42*time.Second, backoff.WithMaxRetries(5))
	for i, d := range backOffDuration() {
		fmt.Println(i, d.String())
	}
	// Output:
	// 0 42s
	// 1 42s
	// 2 42s
	// 3 42s
	// 4 42s
}

func ExampleNewExponentialDuration() {
	backOffDuration := backoff.NewExponentialDuration(
		backoff.WithInitialInterval(10*time.Second),
		backoff.WithRandomizationFactor(0),
		backoff.WithMaxRetries(5),
	)
	for i, d := range backOffDuration() {
		fmt.Println(i, d.String())
	}
	// Output:
	// 0 10s
	// 1 15s
	// 2 22.5s
	// 3 33.75s
	// 4 50.625s
}

func TestNewConstantDuration(t *testing.T) {
	backOffDuration := backoff.NewConstantDuration(42*time.Second, backoff.WithMaxRetries(5))
	a := time.Now()
	got := slices.Collect(values(backOffDuration()))
	if b := time.Since(a); b > time.Millisecond {
		t.Error("It's taking too long")
	}
	want := slices.Repeat([]time.Duration{42 * time.Second}, 5)
	if !slices.Equal(got, want) {
		t.Logf("\n%v\n%v", got, want)
		t.Error("must be equal")
	}
}

func TestNewExponentialDuration(t *testing.T) {
	backOffDuration := backoff.NewExponentialDuration(
		backoff.WithInitialInterval(10*time.Second),
		backoff.WithRandomizationFactor(0),
		backoff.WithMaxRetries(5),
	)
	a := time.Now()
	got := slices.Collect(values(backOffDuration()))
	if b := time.Since(a); b > time.Millisecond {
		panic("It's taking too long")
	}
	want := []time.Duration{
		10 * time.Second,
		15 * time.Second,
		22500 * time.Millisecond,
		33750 * time.Millisecond,
		50625 * time.Millisecond,
	}
	if !slices.Equal(got, want) {
		t.Error("must be equal")
	}
}

func values[K, V any](i iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		next, stop := iter.Pull2(i)
		defer stop()
		for {
			_, v, ok := next()
			if !ok || !yield(v) {
				break
			}
		}
	}
}
