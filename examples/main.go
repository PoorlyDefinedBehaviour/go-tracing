package main

import (
	"context"

	"github.com/poorlydefinedbehaviour/go-tracing"
)

func main() {
	ctx := tracing.WrapContext(context.Background())

	span1 := tracing.Of(ctx).
		Span("Span 1").
		Fields(tracing.Fields{"key_a": 1})

	span2 := tracing.Of(ctx).
		Span("Span 2").
		Fields(tracing.Fields{"key_b": 1})

	tracing.Of(ctx).Info("AAA")

	span1.End()

	tracing.Of(ctx).Info("BBB")

	span2.End()

	// [START - Span 1]
	// [START - Span 2]
	// [INFO - Span 2] AAA - map[key_a:1 key_b:1]
	// [END - Span 1]
	// [INFO - Span 2] BBB - map[key_b:1]
	// [END - Span 2]
}
