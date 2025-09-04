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

package lotel_test

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log/logtest"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/onsi/gomega"

	. "github.com/thediveo/otelcheck/lotel"
)

func ExampleHaveAttribute() {
	/* only in testable example */ Ω := gomega.NewGomega(func(message string, _ ...int) { panic(message) })

	record := logtest.RecordFactory{
		Resource: resource.NewWithAttributes("example.org/foobar",
			attribute.KeyValue{Key: "resource.name", Value: attribute.StringValue("foobar")},
			attribute.KeyValue{Key: "resource.zzz", Value: attribute.IntValue(12345)}),
		InstrumentationScope: &instrumentation.Scope{
			Attributes: attribute.NewSet(attribute.Int("scope.id", 42)),
		},
		Attributes: []log.KeyValue{{Key: "foo", Value: log.StringValue("bar")}},
	}.NewRecord()

	// notice how we easily match attributes on different levels of hierarchy:
	// be it the immediate log record attributes, resource attributes, or
	// instrumentation/scope attributes. It's simply all attributes.
	Ω.Expect(record).To(BeARecord(
		HaveAttribute("resource.name=foobar"),
		HaveAttribute("scope.id"),
		HaveAttribute("foo=bar")))
	// Output:
}
