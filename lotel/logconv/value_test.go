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

package logconv

import (
	"reflect"

	"go.opentelemetry.io/otel/log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/otelcheck/internal/iff"
)

var _ = Describe("log (key-)value conversions", func() {

	DescribeTable("any-fying log values",
		func(v log.Value, expected any) {
			Expect(Any(v)).To(If(expected == nil, BeNil(), Equal(expected)))
		},
		Entry("empty", log.Value{}, nil),
		Entry("bool", log.BoolValue(true), true),
		Entry("int", log.IntValue(42), int64(42)),
		Entry("float", log.Float64Value(42), float64(42)),
		Entry("string", log.StringValue("foo"), "foo"),
		Entry("bytes", log.BytesValue([]byte("bar")), []byte("bar")),
		Entry("slice",
			log.SliceValue(log.StringValue("barz"), log.Int64Value(666)),
			[]any{"barz", int64(666)},
		),
		Entry("map",
			log.MapValue(log.KeyValue{Key: "barz", Value: log.Float64Value(123)},
				log.KeyValue{Key: "foo", Value: log.StringValue("bar")}),
			map[string]any{"foo": "bar", "barz": float64(123)},
		),
	)

	DescribeTable("canonizing any values",
		func(v any, expected any) {
			Expect(reflect.DeepEqual(Canonize(v), expected)).To(BeTrue())
		},
		Entry("nil", nil, nil),
		Entry("bool", true, true),
		Entry("int", 42, int64(42)),
		Entry("int64", int64(42), int64(42)),
		Entry("float32", float32(42), float64(42)),
		Entry("float64", float64(42), float64(42)),
		Entry("string", "foo", "foo"),
		Entry("[]byte", []byte("foo"), []byte("foo")),
		Entry("[]any", []any{"foo", 42, nil}, []any{"foo", int64(42), nil}),
		Entry("map[string]any",
			map[string]any{"foo": "bar", "baz": 42},
			map[string]any{"foo": "bar", "baz": int64(42)}),
	)

	It("panics when canonizing fails", func() {
		Expect(func() {
			_ = Canonize(make(chan struct{}))
		}).To(Panic())
	})

	DescribeTable("any to log value",
		func(v any, expected log.Value) {
			Expect(Value(v).Equal(expected))
		},
		Entry("empty", nil, log.Value{}),
		Entry("bool", true, log.BoolValue(true)),
		Entry("int", 42, log.IntValue(42)),
		Entry("int64", int64(42), log.IntValue(42)),
		Entry("float32", float32(42.0), log.Float64Value(42)),
		Entry("float64", float64(42), log.Float64Value(42)),
		Entry("string", "foo", log.StringValue("foo")),
		Entry("bytes", []byte("bar"), log.BytesValue([]byte("bar"))),
		Entry("slice",
			[]any{"barz", int64(666)},
			log.SliceValue(log.StringValue("barz"), log.Int64Value(666)),
		),
		Entry("map",
			map[string]any{"foo": "bar", "barz": float64(123)},
			log.MapValue(log.KeyValue{Key: "barz", Value: log.Float64Value(123)},
				log.KeyValue{Key: "foo", Value: log.StringValue("bar")}),
		),
	)

	It("panics when an any value cannot be represented as a log.Value", func() {
		Expect(func() {
			_ = Value(new(chan struct{}))
		}).To(PanicWith("logconv.Value: unsupported type *chan struct {}"))
	})

	It("equals", func() {
		mv := Value(map[string]string{
			"foo": "bar",
		})
		Expect(Value(Any(mv)).Equal(mv)).To(BeTrue())
	})

})
