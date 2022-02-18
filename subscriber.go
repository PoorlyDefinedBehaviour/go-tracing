package tracing

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
	// Called when a span is created.
	OnSpanEnter(span string)
	// Called when a span is exited.
	OnSpanExit(span string)
	// Called to log events.
	OnEvent(event Event)
}
