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
	"time"

	"go.opentelemetry.io/otel/sdk/log/logtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaveTimestamp matcher", func() {

	It("matches a plain timestamp", func() {
		t := time.Now()
		r := logtest.RecordFactory{EventName: "foo"}.NewRecord()
		r.SetTimestamp(t)
		Expect(r).To(HaveTimestamp(t))
	})

	It("matches the timestamp using another matcher", func() {
		t := time.Now()
		r := logtest.RecordFactory{EventName: "foo"}.NewRecord()
		r.SetTimestamp(t)
		Expect(r).NotTo(HaveTimestamp(BeTemporally("<", t)))
		Expect(r).To(HaveTimestamp(BeTemporally(">=", t)))
	})

	It("doesn't match a different timestamp", func() {
		r := logtest.RecordFactory{EventName: "foo"}.NewRecord()
		Expect(r).NotTo(HaveTimestamp(time.Now()))
	})

})
