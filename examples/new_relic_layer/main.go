package main

import (
	"context"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/poorlydefinedbehaviour/go-tracing"
	"github.com/poorlydefinedbehaviour/go-tracing/layers"
)

func main() {
	app, _ := newrelic.NewApplication(
		// Name your application
		newrelic.ConfigAppName("Your Application Name"),
		// Fill in your New Relic license key
		newrelic.ConfigLicense("__YOUR_NEW_RELIC_LICENSE_KEY__"),
		// Add logging:
		newrelic.ConfigDebugLogger(os.Stdout),
		// Optional: add additional changes to your configuration via a config function:
		func(cfg *newrelic.Config) {
			cfg.CustomInsightsEvents.Enabled = false
		},
	)

	tracing.WithLayers(layers.NewRelic(app))

	ctx := tracing.WrapContext(context.Background())

	defer tracing.Of(ctx).
		Span("Span 1").
		Fields(tracing.Fields{"key_a": 1}).
		End()

	tracing.Of(ctx).Info("hello world")

	// [START - Span 1]
	// [INFO - Span 1] hello world - map[key_a:1]
	// [END - Span 1]
}
