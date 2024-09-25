package backoff

import (
	"iter"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// BackOffDuration is a function to create an iterator that emits interval durations.
type BackOffDuration func() iter.Seq2[int, time.Duration]

func newBackOffDuration(fn func() backoff.BackOff) BackOffDuration {
	return func() iter.Seq2[int, time.Duration] {
		return func(yield func(int, time.Duration) bool) {
			for i, b := 0, fn(); ; i++ {
				d := b.NextBackOff()
				if d == backoff.Stop || !yield(i, d) {
					break
				}
			}
		}
	}
}

// NewConstantDuration returns a [BackOffDuration], which creates an iterator
// with a backoff policy that always returns the same backoff delay.
//
// If interval < 0 then NewConstantDuration returns the iterator powered by [backoff.StopBackOff].
//
// If interval == 0 then NewConstantDuration returns the iterator powered by [backoff.ZeroBackOff].
//
// If interval > 0 then NewConstantDuration returns the iterator powered by [backoff.ConstantBackOff].
func NewConstantDuration(
	interval time.Duration,
	options ...ConstantOption,
) BackOffDuration {
	co := newConstantOption(interval, options...)
	return newBackOffDuration(co.New)
}

// NewExponentialDuration returns a [BackOffDuration], which creates an iterator powered by [backoff.ExponentialBackOff].
func NewExponentialDuration(
	options ...ExponentialOption,
) BackOffDuration {
	eo := newExponentialOption(options...)
	return newBackOffDuration(eo.New)
}
