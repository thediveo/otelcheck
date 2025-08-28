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

var _ = Describe("HaveBody matcher", func() {

	It("matches a body string value", func() {
		r := logtest.RecordFactory{Body: log.StringValue("doh!")}.NewRecord()
		Expect(r).To(HaveBody("doh!"))
	})

	It("fails to match a body value", func() {
		r := logtest.RecordFactory{Body: log.StringValue("doh!")}.NewRecord()
		Expect(r).NotTo(HaveBody("D'OH!!!"))
	})

	It("matches a body value using another matcher", func() {
		r := logtest.RecordFactory{Body: log.StringValue("doh!")}.NewRecord()
		Expect(r).To(HaveBody(WithTransform(logconv.Any, HavePrefix("do"))))
	})

	It("matches a body value", func() {
		r := logtest.RecordFactory{Body: log.StringValue("doh!")}.NewRecord()
		Expect(r).To(HaveBody(logconv.Value("doh!")))
	})

})
