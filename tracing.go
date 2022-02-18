package tracing

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var subscriber Subscriber = &StdoutSubscriber{}

func SetSubscriber(sub Subscriber) {
	subscriber = sub
}

// Returns a tracing aware context.Context
func WrapContext(ctx context.Context) context.Context {
	return &TracingContext{
		ctx: ctx,
	}
}

type Fields map[string]interface{}

type Tracer struct {
	// Unique id.
	id string
	// Must be locked before accessing any mutable state.
	mu sync.RWMutex
	// The stack of spans we entered by calling Span(string).
	spans []*Spanner
	// Our type wrapping the context that we are allowed to mutate.
	ctx *TracingContext
	// Fields added for context sensitive logging.
	fields Fields
}

func newTracer(ctx *TracingContext) *Tracer {
	return &Tracer{
		id:     uuid.NewString(),
		mu:     sync.RWMutex{},
		spans:  make([]*Spanner, 0),
		fields: make(Fields),
		ctx:    ctx,
	}
}

// Returns a Tracer from the context if it exists,
// otherwise returns a new Tracer.
func Of(ctx context.Context) *Tracer {
	tracingContext, ok := ctx.(*TracingContext)
	if !ok {
		panic("tracing.Of must be called with a tracing context. Use tracing.WrapContext(ctx) to get a TracingContext.")
	}

	if tracer := tracingContext.currentTracer(); tracer != nil {
		return tracer
	}

	return newTracer(tracingContext)
}

// Enters a new span.
func (tracer *Tracer) Span(name string) *Spanner {
	tracer.mu.Lock()
	defer tracer.mu.Unlock()
	spanIndex := len(tracer.spans)
	span := &Spanner{name: name, index: spanIndex, tracer: tracer}
	tracer.spans = append(tracer.spans, span)
	tracer.ctx.addTracerToContext(tracer)
	subscriber.OnSpanEnter(name)
	return span
}

// Merges the key value pairs in `fields` with the tracer fields.
func (tracer *Tracer) Fields(fields Fields) *Tracer {
	for key, value := range fields {
		tracer.fields[key] = value
	}

	return tracer
}

// Called when a span exits.
func (tracer *Tracer) onSpanExit(spanner *Spanner) {
	tracer.mu.Lock()
	defer tracer.mu.Unlock()

	// Let the subscriber know the span exited
	subscriber.OnSpanExit(spanner.name)

	// If the last span exited, we can make this fast
	// by just removing the last element from the slice.
	// This will be the hot path.
	if spanner.index == len(tracer.spans)-1 {
		tracer.spans = tracer.spans[:len(tracer.spans)-1]
		return
	}

	// If a span that's in the middle of the list just exited,
	// we take the slow path.
	tracer.spans = append(tracer.spans[:spanner.index], tracer.spans[spanner.index+1:]...)
}

// Returns the newest span.
func (tracer *Tracer) currentSpan() *Spanner {
	tracer.mu.RLock()
	defer tracer.mu.RUnlock()

	numSpans := len(tracer.spans)
	if numSpans == 0 {
		return nil
	}

	return tracer.spans[numSpans-1]
}

func (tracer *Tracer) sendEventToSubscriber(level EventLevel, format string, values ...interface{}) {
	event := Event{
		Level:   level,
		Fields:  tracer.fields,
		Span:    tracer.currentSpan(),
		Message: fmt.Sprintf(format, values...),
	}

	subscriber.OnEvent(event)
}

func (tracer *Tracer) Info(format string, values ...interface{}) {
	tracer.sendEventToSubscriber(INFO, format, values...)
}

func (tracer *Tracer) Warn(format string, values ...interface{}) {
	tracer.sendEventToSubscriber(WARN, format, values...)
}

func (tracer *Tracer) Error(format string, values ...interface{}) {
	tracer.sendEventToSubscriber(ERROR, format, values...)
}

func (tracer *Tracer) Debug(format string, values ...interface{}) {
	tracer.sendEventToSubscriber(DEBUG, format, values...)
}
