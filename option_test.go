package backoff

import (
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
)

func TestNewConstantOption(t *testing.T) {
	opt := newConstantOption(
		1234,
		WithMaxRetries(7890),
	)
	if opt.Interval != 1234 {
		t.Error("interval")
	}
	if opt.MaxRetries != 7890 {
		t.Error("WithMaxRetries")
	}
}

func TestNewExponentialOption(t *testing.T) {
	opt := newExponentialOption(
		WithInitialInterval(1234),
		WithRandomizationFactor(2345),
		WithMultiplier(3456),
		WithMaxInterval(4567),
		WithMaxElapsedTime(5678),
		WithRetryStopDuration(6789),
		WithClockProvider(theClock),
	)
	b := opt.New().(*backoff.ExponentialBackOff)

	if b.InitialInterval != 1234 {
		t.Error("WithInitialInterval")
	}
	if b.RandomizationFactor != 2345 {
		t.Error("WithRandomizationFactor")
	}
	if b.Multiplier != 3456 {
		t.Error("WithMultiplier")
	}
	if b.MaxInterval != 4567 {
		t.Error("WithMaxInterval")
	}
	if b.MaxElapsedTime != 5678 {
		t.Error("WithMaxElapsedTime")
	}
	if b.Stop != 6789 {
		t.Error("WithStop")
	}
	if b.Clock != theClock {
		t.Errorf("WithClock %#v %#v", b.Clock, theClock)
	}
	if opt.MaxRetries != 0 {
		t.Error("WithMaxRetries")
	}

	opt = newExponentialOption(
		WithMaxRetries(7890),
	)
	if opt.MaxRetries != 7890 {
		t.Error("WithMaxRetries")
	}
}

var theClock testClock

type testClock struct{}

func (testClock) Now() time.Time { return time.Now() }
