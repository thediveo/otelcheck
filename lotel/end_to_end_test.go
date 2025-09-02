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
	"context"

	"github.com/thediveo/otelcheck/lotel/testlogger"
	"go.opentelemetry.io/otel/log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("lotel e2e", func() {

	It("asserts log records", func(ctx context.Context) {
		l, ch := testlogger.New(10)

		go func() {
			defer GinkgoRecover()
			r := log.Record{}
			r.SetEventName("org.foo")
			r.AddAttributes(log.KeyValue{Key: "foo", Value: log.IntValue(42)},
				log.KeyValue{Key: "bar", Value: log.StringValue("barf!")})
			l.Emit(ctx, r)
		}()

		Eventually(ch).Should(
			Receive(BeARecord(
				HaveEventName("org.foo"),
				HaveAttributeWithValue("foo", 42),
				HaveAttribute("bar=barf!"),
			)))

	})

})
