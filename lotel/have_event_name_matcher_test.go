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
	"go.opentelemetry.io/otel/sdk/log/logtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaveEventName matcher", func() {

	It("matches a plain string name", func() {
		r := logtest.RecordFactory{EventName: "foo"}.NewRecord()
		Expect(r).To(HaveEventName("foo"))
	})

	It("matches a matcher", func() {
		r := logtest.RecordFactory{EventName: "foo"}.NewRecord()
		Expect(r).To(HaveEventName(HavePrefix("fo")))
	})

	It("doesn't match a different plain string name", func() {
		r := logtest.RecordFactory{EventName: "foo"}.NewRecord()
		Expect(r).NotTo(HaveEventName("bar"))
	})

})
