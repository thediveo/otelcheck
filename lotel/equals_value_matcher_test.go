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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EqualsValue matcher", func() {

	DescribeTable("log.Value matching success",
		func(actual, expected log.Value) {
			Expect(EqualsValue(expected).Match(actual)).To(BeTrue())
		},
		Entry(nil, log.StringValue("foo"), log.StringValue("foo")),
		Entry(nil, log.IntValue(42), log.Int64Value(42)),
		Entry(nil,
			log.SliceValue(log.StringValue("bar"), log.IntValue(42)),
			log.SliceValue(log.StringValue("bar"), log.Int64Value(42))),
	)

	DescribeTable("log.Value matching fails",
		func(actual, expected log.Value) {
			Expect(EqualsValue(expected).Match(actual)).To(BeFalse())
		},
		Entry(nil, log.StringValue("foo"), log.StringValue("bar")),
		Entry(nil, log.IntValue(42), log.Float64Value(42)),
		Entry(nil,
			log.SliceValue(log.StringValue("bar"), log.IntValue(42)),
			log.SliceValue(log.StringValue("barz"), log.Int64Value(42))),
	)

})
