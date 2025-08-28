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
	"github.com/thediveo/otelcheck/lotel/logconv"
)

var _ = Describe("BeARecord matcher", func() {

	It("fails when not given a log record", func() {
		m := BeARecord(nil)
		Expect(m.Match(42)).Error().To(HaveOccurred())
	})

	It("matches when not all matchers are satisfied", func() {
		m := BeARecord(HaveField("Body().AsString()", "doh!"), HaveField("EventName()", "uh-oh"))
		r := logtest.RecordFactory{Body: log.StringValue("doh!")}.NewRecord()
		r.SetEventName("uh-oh")
		Expect(m.Match(r)).To(BeTrue())
	})

	It("doesn't match when not all matchers are satisfied", func() {
		m := BeARecord(HaveField("Body().String()", "doh!"), HaveField("EventName()", "bar"))
		r := logtest.RecordFactory{Body: log.StringValue("doh!")}.NewRecord()
		Expect(m.Match(r)).To(BeFalse())

		r.SetBody(log.BoolValue(false))
		Expect(m.Match(r)).To(BeFalse())
	})

	It("matches attributes", func() {
		r := logtest.RecordFactory{
			Attributes: []log.KeyValue{{Key: "foo", Value: logconv.Value("bar")}},
		}.NewRecord()
		Expect(r).To(BeARecord(HaveAttribute("foo")))
		Expect(r).To(BeARecord(HaveAttribute("foo=bar")))
		Expect(r).NotTo(BeARecord(HaveAttribute("foo"), HaveAttribute("bar")))
	})

})
