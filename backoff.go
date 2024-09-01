// Package backoff provides iterators powered by github.com/cenkalti/backoff/v4.
//
// This package provides the type [BackOff], which creates an iterator powered by
// github.com/cenkalti/backoff/v4.
//
// [BackOff] can be created by [NewConstant] and [NewExponential].
//
// There are four policies.
//
//   - [backoff.StopBackOff]
//   - [backoff.ZeroBackOff]
//   - [backoff.ConstantBackOff]
//   - [backoff.ExponentialBackOff]
//
// Iteratros powered by StopBackOff, ZeroBackOff and ConstantBackOff can be created from [NewConstant].
//
// Iterators powered by ExponentialBackOff can be created from [NewExponential].
//
// Iterators are used in `for statements with range clause` so that iteration can be easily terminated with `break`.
// Alternatively, you can use [WithMaxRetries] to limit the number of retries.
//
// You can interrupt the wait by canceling the context you specified when creating the iterator, idiomatically.
//
// # See
//
//   - https://go.dev/ref/spec#For_range
//   - https://go.dev/doc/go1.23#iterators
package backoff

import (
	"context"
	"errors"
	"iter"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// BackOff is a function to create an iterator.
// The iterator can be cancelled by a context.
// It is safe to pass nil for the context when creating an iterator.
// The iterator yields sequence of integers representing the number of retry attempts.
type BackOff func(context.Context) iter.Seq[int]

var errContinue = errors.New("continue")

func newBackOff(fn func() backoff.BackOff) BackOff {
	return func(ctx context.Context) iter.Seq[int] {
		return func(yield func(int) bool) {
			i := 0
			task := func() error {
				if !yield(i) {
					return nil // stops the retry loop.
				}
				i++
				return errContinue // continue the retry loop
			}
			b := fn()
			if ctx != nil {
				b = backoff.WithContext(b, ctx)
			}
			backoff.Retry(task, b)
		}
	}
}

// NewConstant returns a [BackOff], which creates an iterator with a backoff policy that
// always returns the same backoff delay.
//
// If interval < 0 then NewConstant returns the iterator powered by [backoff.StopBackOff].
//
// If interval == 0 then NewConstant returns the iterator powered by [backoff.ZeroBackOff].
//
// If interval > 0 then NewConstant returns the iterator powered by [backoff.ConstantBackOff].
func NewConstant(interval time.Duration, options ...ConstantOption) BackOff {
	co := newConstantOption(interval, options...)
	return newBackOff(co.New)
}

// NewExponential returns a [BackOff], which creates an iterator powered by [backoff.ExponentialBackOff].
func NewExponential(options ...ExponentialOption) BackOff {
	eo := newExponentialOption(options...)
	return newBackOff(eo.New)
}

func WithMaxRetries(max uint64) Option {
	return maxRetries(max)
}

type maxRetries uint64

func (o maxRetries) applyToConstantOption(co *constantOption) {
	co.MaxRetries = uint64(o)
}

func (o maxRetries) applyToExponentialOption(eo *exponentialOption) {
	eo.MaxRetries = uint64(o)
}

// Option is a optional parameter for [NewConstant] and [NewExponential].
type Option interface {
	ConstantOption
	ExponentialOption
}

// ConstantOption is a optional parameter for [NewConstant].
type ConstantOption interface {
	applyToConstantOption(*constantOption)
}

// ExponentialOption is a optional parameter for [NewExponential].
//
// The following options are available:
//
//   - [WithInitialInterval]: The initial interval for the first retry
//   - [WithRandomizationFactor]: The randomization factor to use for creating a range around the retry interval
//   - [WithMultiplier]: The factor by which the retry interval increases
//   - [WithMaxInterval]: The maximum value of the retry interval
//   - [WithMaxElapsedTime]: The maximum amount of time to retry
//   - [WithRetryStopDuration]: The interval at which the backoff stops increasing
//   - [WithClockProvider]: The clock to use for time measurements
//
// See [backoff.ExponentialBackOff] for more details on each option.
type ExponentialOption interface {
	applyToExponentialOption(*exponentialOption)
}

type constantOption struct {
	MaxRetries uint64
	Interval   time.Duration
}

func newConstantOption(interval time.Duration, options ...ConstantOption) *constantOption {
	co := &constantOption{
		Interval: interval,
	}
	for _, o := range options {
		o.applyToConstantOption(co)
	}
	return co
}

func (co *constantOption) New() backoff.BackOff {
	if co.Interval < 0 {
		return &backoff.StopBackOff{}
	}
	var b backoff.BackOff
	if co.Interval == 0 {
		b = &backoff.ZeroBackOff{}
	} else {
		b = backoff.NewConstantBackOff(co.Interval)
	}
	if co.MaxRetries > 0 {
		b = backoff.WithMaxRetries(b, co.MaxRetries)
	}
	return b
}

type exponentialOption struct {
	Options    []backoff.ExponentialBackOffOpts
	MaxRetries uint64
}

func newExponentialOption(options ...ExponentialOption) *exponentialOption {
	eo := &exponentialOption{}
	for _, o := range options {
		o.applyToExponentialOption(eo)
	}
	return eo
}

func (eo *exponentialOption) New() (b backoff.BackOff) {
	b = backoff.NewExponentialBackOff(eo.Options...)
	if eo.MaxRetries > 0 {
		b = backoff.WithMaxRetries(b, eo.MaxRetries)
	}
	return b
}

type opts backoff.ExponentialBackOffOpts

var _ ExponentialOption = (*opts)(nil)

func newOpts(o backoff.ExponentialBackOffOpts) ExponentialOption {
	return opts(o)
}

func (o opts) applyToExponentialOption(eo *exponentialOption) {
	eo.Options = append(eo.Options, backoff.ExponentialBackOffOpts(o))
}

// WithInitialInterval sets the initial interval between retries.
//
// See [backoff.WithInitialInterval].
func WithInitialInterval(d time.Duration) ExponentialOption {
	return newOpts(backoff.WithInitialInterval(d))
}

// WithRandomizationFactor sets the randomization factor to add jitter to intervals.
//
// See [backoff.WithRandomizationFactor].
func WithRandomizationFactor(v float64) ExponentialOption {
	return newOpts(backoff.WithRandomizationFactor(v))
}

// WithMultiplier sets the multiplier for increasing the interval after each retry.
//
// See [backoff.WithMultiplier].
func WithMultiplier(v float64) ExponentialOption {
	return newOpts(backoff.WithMultiplier(v))
}

// WithMaxInterval sets the maximum interval between retries.
//
// See [backoff.WithMaxInterval].
func WithMaxInterval(d time.Duration) ExponentialOption {
	return newOpts(backoff.WithMaxInterval(d))
}

// WithMaxElapsedTime sets the maximum total time for retries.
//
// See [backoff.WithMaxElapsedTime].
func WithMaxElapsedTime(d time.Duration) ExponentialOption {
	return newOpts(backoff.WithMaxElapsedTime(d))
}

// WithRetryStopDuration sets the duration after which retries should stop.
//
// See [backoff.WithRetryStopDuration].
func WithRetryStopDuration(d time.Duration) ExponentialOption {
	return newOpts(backoff.WithRetryStopDuration(d))
}

// WithClockProvider sets the clock used to measure time.
//
// See [backoff.WithClockProvider].
func WithClockProvider(c backoff.Clock) ExponentialOption {
	return newOpts(backoff.WithClockProvider(c))
}
