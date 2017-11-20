package sinks

import (
	"context"

	"github.com/stripe/veneur/samplers"
	"github.com/stripe/veneur/ssf"
	"github.com/stripe/veneur/trace"
)

// MetricSink is a receiver of `InterMetric`s when Veneur periodically flushes
// it's aggregated metrics.
type MetricSink interface {
	Name() string
	// Start finishes setting up the sink and starts any
	// background processing tasks that the sink might have to run
	// in the background. It's invoked when the server starts.
	Start(traceClient *trace.Client) error
	// Flush receives `InterMetric`s from Veneur and is responsible for "sinking"
	// these metrics to whatever it's backend wants. Note that the sink must
	// **not** mutate the incoming metrics as they are shared with other sinks.
	Flush(context.Context, []samplers.InterMetric) error
	// This one is temporary?
	FlushEventsChecks(ctx context.Context, events []samplers.UDPEvent, checks []samplers.UDPServiceCheck)
}

// SpanSink is a receiver of spans that handles sending those spans to some
// downstream sink. Calls to `Ingest(span)` are meant to give the sink control
// of the span, with periodic calls to flush as a signal for sinks that don't
// handle their own flushing in a separate goroutine, etc. Note that SpanSinks
// differ from MetricSinks because Veneur does *not* aggregate Spans.
type SpanSink interface {
	// Start finishes setting up the sink and starts any
	// background processing tasks that the sink might have to run
	// in the background. It's invoked when the server starts.
	Start(*trace.Client) error
	Name() string
	// Flush receives `SSFSpan`s from Veneur **as they arrive**. If the sink wants
	// to buffer spans it may do so and defer sending until `Flush` is called.
	Ingest(ssf.SSFSpan) error
	// Invoked at the same interval as metric flushes, this can be used as a
	// signal for the sink to write out if it was buffering or something.
	Flush()
}