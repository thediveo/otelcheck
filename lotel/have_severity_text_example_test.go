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
	"go.opentelemetry.io/otel/sdk/log/logtest"

	"github.com/onsi/gomega"

	. "github.com/thediveo/otelcheck/lotel"
)

func ExampleHaveSeverityText() {
	/* only in testable example */ Ω := gomega.NewGomega(func(message string, _ ...int) { panic(message) })

	record := logtest.RecordFactory{SeverityText: "foobar"}.NewRecord()

	Ω.Expect(record).To(HaveSeverityText("foobar"))
	Ω.Expect(record).To(HaveSeverityText(gomega.HavePrefix("foo")))
	// Output:
}
