package layers

import (
	"errors"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/poorlydefinedbehaviour/go-tracing"
)

type NewRelicApplication interface {
	StartTransaction(name string) NewRelicTransaction
}

type NewRelicTransaction interface {
	AddAttribute(key string, value interface{})
	NoticeError(err error)
	End()
}

type newRelicLayer struct {
	app NewRelicApplication
	// Maps a span to its new relic transaction.
	transactions map[string]NewRelicTransaction
}

type newRelicAdapter struct {
	app *newrelic.Application
}

func (adapter *newRelicAdapter) StartTransaction(name string) NewRelicTransaction {
	trx := adapter.app.StartTransaction(name)
	return &newRelicTransactionAdapter{trx: trx}
}

type newRelicTransactionAdapter struct {
	trx *newrelic.Transaction
}

func (adapter *newRelicTransactionAdapter) AddAttribute(key string, value interface{}) {
	adapter.trx.AddAttribute(key, value)
}

func (adapter *newRelicTransactionAdapter) NoticeError(err error) {
	adapter.trx.NoticeError(err)
}

func (adapter *newRelicTransactionAdapter) End() {
	adapter.trx.End()
}

func NewRelic(app *newrelic.Application) tracing.Layer {
	return &newRelicLayer{
		app:          &newRelicAdapter{app: app},
		transactions: make(map[string]NewRelicTransaction),
	}
}

// Called when a span is created.
func (layer *newRelicLayer) OnSpanEnter(span string) {
	layer.transactions[span] = layer.app.StartTransaction(span)
}

func (layer *newRelicLayer) OnSpanExit(span string) {
	trx := layer.transactions[span]
	if trx == nil {
		return
	}
	trx.End()
	delete(layer.transactions, span)
}

func (layer *newRelicLayer) OnEvent(event *tracing.Event) {
	if event.Span == nil {
		return
	}

	trx := layer.transactions[event.Span.Name]
	if trx == nil {
		return
	}

	for key, value := range event.Fields {
		trx.AddAttribute(key, value)
	}

	if event.Level == tracing.ERROR {
		trx.NoticeError(errors.New(event.Message))
	}
}
