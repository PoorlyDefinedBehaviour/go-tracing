package tracing

type Layer interface {
	// Called when a span is created.
	OnSpanEnter(span string)
	// Called when a span is exited.
	OnSpanExit(span string)
	// Called to log events.
	OnEvent(event Event)
}
