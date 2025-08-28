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

package testlogger

import (
	"context"

	"github.com/thediveo/otelcheck/exporters/chanlog"
	"go.opentelemetry.io/otel/log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OTel test logger", func() {

	It("automatically shuts down the processor and exporter, closing the record channel", func(ctx context.Context) {
		var ch chanlog.RecordsChannel
		DeferCleanup(func() {
			By("end of test")
			Expect(ch).To(HaveLen(1))
			// drain the canary log record as otherwise BeClose() will fail:
			// while the channel actually is closed, it would still contain
			// element(s).
			<-ch
			Expect(ch).To(BeClosed())
		})
		l, ch := New(42)
		r := log.Record{}
		r.SetBody(log.StringValue("bah!"))
		l.Emit(ctx, r)
	})

})
