package tracing

import "fmt"

type EventLevel struct {
	value string
}

func (level *EventLevel) String() string {
	return level.value
}

var (
	INFO  = EventLevel{value: "INFO"}
	WARN  = EventLevel{value: "WARN"}
	ERROR = EventLevel{value: "ERROR"}
	DEBUG = EventLevel{value: "DEBUG"}
)

type Event struct {
	Level   EventLevel
	Fields  Fields
	Span    *Spanner
	Message string
}

type Subscriber interface {
	// Adds a layer that is called before events are logged.
	// Layers can return a modified events.
	With(layer Layer)
	// Called when a span is created.
	OnSpanEnter(span string)
	// Called when a span is exited.
	OnSpanExit(span string)
	// Called to log events.
	OnEvent(event Event)
}

type StdoutSubscriber struct{}

func (subscriber *StdoutSubscriber) OnSpanEnter(span string) {
	fmt.Printf("[START - %s]\n", span)
}

func (subscriber *StdoutSubscriber) OnSpanExit(span string) {
	fmt.Printf("[END - %s]\n", span)
}

func (subscriber *StdoutSubscriber) OnEvent(event Event) {
	if event.Span != nil {
		fmt.Printf("[%s - %s] %s - %+v\n", event.Level.String(), event.Span.name, event.Message, event.Fields)
	} else {
		fmt.Printf("[%s] %s - %+v\n", event.Level.String(), event.Message, event.Fields)
	}
}

func (subscriber *StdoutSubscriber) With(layer Layer) {
	panic("TODO")
}
