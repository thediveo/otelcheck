package chanlog_test

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"

	"github.com/thediveo/otelcheck/exporters/chanlog"
)

// Please consider using [github.com/thediveo/otelcheck/lotel/testlogger.New]
// instead.
func Example() {
	// create a test logger that emits log records into a buffered channel with
	// a capacity for 10 log records, the detailed way. Consider using
	// testlogger.New instead.
	exporter, _ := chanlog.New(chanlog.WithCap(10))
	processor := sdklog.NewSimpleProcessor(exporter)
	provider := sdklog.NewLoggerProvider(sdklog.WithProcessor(processor))
	logger := provider.Logger("testlogger")
	defer func() { _ = processor.Shutdown(context.TODO()) }()

	// get the channel the records are sent to now, as the exporter's Ch() will return a nil
	// channel as soon as Shutdown has been called on the exporter.
	ch := exporter.Ch()

	r := log.Record{}
	r.SetBody(log.StringValue("DO'H!"))
	logger.Emit(context.TODO(), r)

	select {
	case r := <-ch:
		fmt.Println(r.Body().AsString())
	case <-time.After(5 * time.Second):
		panic("expected to receive a log record")
	}
	// Output: DO'H!
}
