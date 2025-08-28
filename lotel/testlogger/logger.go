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
	"github.com/thediveo/otelcheck/exporters/chanlog"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"

	gi "github.com/onsi/ginkgo/v2"
)

// New returns a new OTel logger object together with the log record channel the
// logger exports to. Behind the scenes, New creates a logger provider and
// schedules it do be automatically shut down when the calling specification
// terminates; New then wires up the logger provider to an exporter (using an
// [sdklog.SimpleProcessor]) that feeds into a Go chan buffering log records.
//
// # Notes
//
// In the OTel SDK, individual [log.Logger] objects cannot and don't need to be
// shut down. However, the logger provider together with its processor(s) and
// exporter(s) need to be shut down, using [sdklog.LoggerProvider.Shutdown]. For
// convenience, we don't expose the throw-away logger provider and thus
// automatically schedule for a deferred shutdown cleanup via [gi.DeferCleanup].
func New(capacity int) (log.Logger, chanlog.RecordsChannel) {
	gi.GinkgoHelper()

	exp, _ := chanlog.New(chanlog.WithCap(capacity))
	proc := sdklog.NewSimpleProcessor(exp)
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(proc))
	gi.DeferCleanup(lp.Shutdown)
	l := lp.Logger("testlogger")

	return l, exp.Ch()
}
