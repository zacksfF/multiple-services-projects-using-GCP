// PKG metrics constructs the metrics the application will track.
package metrics

import (
	"context"
	"expvar"
	"runtime"
)

// This holds the single instance of the metrics value needed for
// collecting metrics.
var met metrics

// metrics represent the set of metrics we gather.
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

// init costruct teh metrics value taht will be used to capture metrics
func init() {
	met = metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("request"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

type ctxKey int

const key ctxKey = 1

// Set sets the metrics data into the context.
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, &met)
}

// AddGoRoutines refreshes the goroutines metric.
func AddGoRoutines(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		g := int64(runtime.NumGoroutine())
		v.goroutines.Set(g)
		return g
	}
	return 0
}

// AddRequest increment the request metric by 1.
func AddRequest(ctx context.Context) int64 {
	v, ok := ctx.Value(key).(*metrics)
	if ok {
		v.requests.Add(1)
		return v.requests.Value()
	}
	return 0
}

// AddErrors increments the errors metric by 1.
func AddErrors(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
		return v.errors.Value()
	}

	return 0
}

// AddPanics increments the panics metric by 1.
func AddPanics(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
		return v.panics.Value()
	}

	return 0
}
