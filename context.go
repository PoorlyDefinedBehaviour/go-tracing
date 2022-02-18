package tracing

import (
	"context"
	"time"
)

// context.Context compatible type that provides
// tracing support.
type TracingContext struct {
	ctx context.Context
}

var contextTracerKey = &struct {
	value string
}{
	value: "contextTracerKey",
}

func (tracingCtx *TracingContext) addTracerToContext(tracer *Tracer) {
	tracingCtx.ctx = context.WithValue(tracingCtx.ctx, contextTracerKey, tracer)
}

// Returns that tracer being held by the internal context, if any.
func (tracingCtx *TracingContext) currentTracer() *Tracer {
	tracer, ok := tracingCtx.ctx.Value(contextTracerKey).(*Tracer)
	if !ok {
		return nil
	}

	return tracer
}

func (tracingCtx *TracingContext) Deadline() (deadline time.Time, ok bool) {
	return tracingCtx.ctx.Deadline()
}

func (tracingCtx *TracingContext) Done() <-chan struct{} {
	return tracingCtx.ctx.Done()
}

func (tracingCtx *TracingContext) Err() error {
	return tracingCtx.ctx.Err()
}

func (tracingCtx *TracingContext) Value(key interface{}) interface{} {
	return tracingCtx.ctx.Value(key)
}
