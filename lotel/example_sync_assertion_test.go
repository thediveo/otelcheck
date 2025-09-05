// Copyright 2025 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package lotel

import (
	"context"
	"slices"
	"time"

	"github.com/thediveo/otelcheck/lotel/testlogger"
	"github.com/thediveo/otelcheck/x/chans"
	"go.opentelemetry.io/otel/log"

	"github.com/onsi/gomega"
)

// This is a complete example showing how to create a logger for testing
// purposes, logging some records, collecting all log records, and then only
// finally asserting the correct records have been received from the channel the
// logger exports into.
func Example_synchronous_assertions() {
	/* only in testable example */ Ω := gomega.NewGomega(func(message string, _ ...int) { panic(message) })

	// only in testable example,so when no suitable test context is at hand.
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	// let's start with the actual example test code now...
	logger, shutdown, ch := testlogger.New(10)
	defer shutdown(ctx)

	// let's log a few records asynchronously.
	go func() {
		r := log.Record{}
		r.SetEventName("org.foo")
		r.AddAttributes(log.KeyValue{Key: "foo", Value: log.IntValue(42)},
			log.KeyValue{Key: "bar", Value: log.StringValue("barf!")})
		logger.Emit(ctx, r)

		r = log.Record{}
		r.SetEventName("org.bar")
		logger.Emit(ctx, r)

		shutdown(ctx)
	}()

	// Wait for all records to arrive and collect them until the channel closes.
	recs := slices.Collect(chans.All(ctx, ch))

	// assert the log records...
	Ω.Expect(recs).To(gomega.ContainElements(
		BeARecord(
			HaveEventName("org.foo"),
			// assert not only log record attributes...
			HaveAttributeWithValue("foo", 42),
			HaveAttribute("bar=barf!"),
			// ...but also resource and instrument/scope attributes.
			HaveAttributeWithValue(testlogger.InstrumentationAttributeName, testlogger.InstrumentationAttributeValue),
			HaveAttributeWithValue("service.name", gomega.HavePrefix("unknown_service:")),
		),
		BeARecord(HaveEventName("org.bar")),
	))

	Ω.Expect(recs).NotTo(
		gomega.ContainElement(BeARecord(HaveEventName("org.barz"))))

	// Output:
}
