package main

import (
	"context"

	"github.com/poorlydefinedbehaviour/go-tracing"
)

func main() {
	ctx := tracing.WrapContext(context.Background())

	defer tracing.Of(ctx).
		Span("My Span Name").
		Fields(tracing.Fields{"request_id": 1}).
		End()

	tracing.Of(ctx).Fields(tracing.Fields{"hello": "world"}).Info("hello world")
	// [START - My Span Name]
	// [INFO - My Span Name] hello world - map[hello:world request_id:1]
	// [END - My Span Name]
}
