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
	"go.opentelemetry.io/otel/sdk/log/logtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/otelcheck/internal/iff"
)

var _ = Describe("HaveSeverity matcher", func() {

	DescribeTable("matching severity levels, or not",
		func(severity log.Severity, expected any, matches bool) {
			r := logtest.RecordFactory{EventName: "foo"}.NewRecord()
			r.SetSeverity(severity)
			If(matches, Assertion.To, Assertion.NotTo)(Expect(r), HaveSeverity(expected))
		},
		Entry(nil, log.SeverityError, log.SeverityError, true),
		Entry(nil, log.SeverityError, log.SeverityDebug, false),
		Entry(nil, log.SeverityError, Not(BeZero()), true),
		Entry(nil, log.SeverityError, BeNumerically(">", log.SeverityDebug), true),
		Entry(nil, log.SeverityError, BeNumerically(">", log.SeverityFatal), false),
	)

})
