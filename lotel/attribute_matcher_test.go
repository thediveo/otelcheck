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
	"go.opentelemetry.io/otel/log"

	"github.com/thediveo/otelcheck/lotel/logconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

})
