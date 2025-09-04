package testlogger_test

import (
	"context"
	"fmt"
	"time"

	"github.com/thediveo/otelcheck/lotel/testlogger"
	"go.opentelemetry.io/otel/log"
)

func Example() {
	// create a test logger that emits log records into a buffered channel with
	// a capacity for 10 log records.
	logger, shutdown, ch := testlogger.New(10)
	defer shutdown(context.TODO())

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
