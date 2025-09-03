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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log/logtest"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/thediveo/otelcheck/lotel/logconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ty "github.com/onsi/gomega/types"
	. "github.com/thediveo/otelcheck/internal/iff"
)

var _ = Describe("HaveAttribute(WithValue) matchers", func() {

	DescribeTable("matches attributes using name[=value] format",
		func(attrspec string, attrv any, match bool) {
			attr := log.KeyValue{Key: "foo", Value: logconv.Value(attrv)}
			If(match, Assertion.To, Assertion.NotTo)(Expect(attr),
				HaveAttribute(attrspec))
		},
		Entry(nil, "foo", "bar", true),
		Entry(nil, "foo=bar", "bar", true),
		Entry(nil, "foo=foo", "bar", false),
		Entry(nil, "foo=", "", true),
		Entry(nil, "foo=", "bar", false),
	)

	DescribeTable("matches attributes using explicit name, value",
		func(name, value any, attrv any, match bool) {
			attr := log.KeyValue{Key: "foo", Value: logconv.Value(attrv)}
			If(match, Assertion.To, Assertion.NotTo)(Expect(attr),
				HaveAttributeWithValue(name, value))
		},
		Entry(nil, "foo", "bar", "bar", true),
		Entry(nil, "foo", "", "bar", false),
		Entry(nil, "foo", "bar", nil, false),
		Entry(nil, "foo", nil, nil, true),
		Entry(nil, "foo", "bar", 42, false),
		Entry(nil, "foo", 42, 42, true),
		Entry(nil, "foo", []any{"untruth", float32(123.0)}, []any{"untruth", float64(123.0)}, true),
	)

	DescribeTable("matches attributes in record, resource, scope",
		func(m ty.GomegaMatcher, match bool) {
			r := logtest.RecordFactory{
				Resource: resource.NewWithAttributes("example.org/foobar",
					attribute.KeyValue{Key: "resource.name", Value: attribute.StringValue("foobar")},
					attribute.KeyValue{Key: "resource.zzz", Value: attribute.IntValue(12345)}),
				InstrumentationScope: &instrumentation.Scope{
					Attributes: attribute.NewSet(attribute.Int("scope.id", 42)),
				},
			}.NewRecord()
			r.AddAttributes(log.String("foo", "bar"))
			If(match, Assertion.To, Assertion.NotTo)(Expect(r), m)
		},
		Entry(nil, HaveAttribute("resource.name=foobar"), true),
		Entry(nil, HaveAttribute("scope.id"), true),
		Entry(nil, HaveAttributeWithValue("scope.id", 42), true),
	)

	DescribeTable("attribute matching error handling",
		func(m ty.GomegaMatcher) {
			r := logtest.RecordFactory{}.NewRecord()
			r.AddAttributes(log.String("foo", "bar"))
			Expect(m.Match(r)).Error().To(HaveOccurred())
		},
		Entry(nil, HaveAttribute(BeTrue())),
	)

	It("returns matching errors when trying to match resource and scope attributes", func() {
		r := logtest.RecordFactory{
			Resource: resource.NewWithAttributes("example.org/foobar",
				attribute.KeyValue{Key: "resource.name", Value: attribute.StringValue("foobar")}),
		}.NewRecord()
		Expect(HaveAttribute(BeTrue()).Match(r)).Error().To(HaveOccurred())

		r = logtest.RecordFactory{
			InstrumentationScope: &instrumentation.Scope{
				Attributes: attribute.NewSet(attribute.Int("scope.id", 42)),
			},
		}.NewRecord()
		Expect(HaveAttribute(BeTrue()).Match(r)).Error().To(HaveOccurred())
	})

})
