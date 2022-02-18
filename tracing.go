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
	span := &Spanner{name: name, tracer: tracer}
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

// Removes `fields` from the tracer fields.
func (tracer *Tracer) removeFields(fields Fields) *Tracer {
	for key := range fields {
		delete(tracer.fields, key)
	}

	return tracer
}

// Called when a span exits.
func (tracer *Tracer) onSpanExit(spanner *Spanner) {
	tracer.mu.Lock()
	defer tracer.mu.Unlock()

	// The fields that belong to the span that is exiting won't be used anymore.
	tracer.removeFields(spanner.fields)

	// If the tracer has a previous span.
	if len(tracer.spans) > 1 {
		previousSpan := tracer.spans[len(tracer.spans)-2]
		// Add the previous span fields back.
		tracer.Fields(previousSpan.fields)
	}

	// Let the subscriber know the span exited
	subscriber.OnSpanExit(spanner.name)

	for i := len(tracer.spans) - 1; i > -1; i-- {
		// If we found the span we want to remove.
		if tracer.spans[i] == spanner {
			// Remove it from the list of spans.
			tracer.spans = append(tracer.spans[:i], tracer.spans[i+1:]...)
			break
		}
	}
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
