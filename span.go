package tracing

type Spanner struct {
	name string
	// The index of the spanner in `Tracer.spans`
	index int
	// The `Tracer` that created this span.
	tracer *Tracer
	fields Fields
}

// Adds the pairs in `keyValuePairs` to the events emitted by the tracer.
func (spanner *Spanner) Fields(fields Fields) *Spanner {
	// Add fields to the tracer
	spanner.tracer.Fields(fields)

	return &Spanner{
		name:   spanner.name,
		index:  spanner.index,
		tracer: *&spanner.tracer,
		fields: fields,
	}
}

func (spanner *Spanner) End() {
	spanner.tracer.onSpanExit(spanner)
}