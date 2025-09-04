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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

const (
	InstrumentationAttributeName  = "github.com/thediveo/lotel/testlogger"
	InstrumentationAttributeValue = 42
)

// New returns a new OTel logger object together with the log record channel the
// logger exports to. Behind the scenes, New creates a logger provider and wires
// it up to an exporter (using an [sdklog.SimpleProcessor]) that feeds into a Go
// chan buffering log records.
//
// Callers should use the returned shutdown function to shut down the logger
// provider as well as the exporter, with the channel also getting closed in the
// process.
//
// The returned logger object features a “canary” instrumentation/scope
// attribute that can be used in testing. This instrument attribute has the name
// [InstrumentationAttributeName] and value [InstrumentationAttributeValue].
//
// # Notes
//
// In the OTel SDK, individual [log.Logger] objects cannot and don't need to be
// shut down. However, the logger provider together with its processor(s) and
// exporter(s) need to be shut down, using [sdklog.LoggerProvider.Shutdown]. We
// don't expose the throw-away logger provider but instead expose an omnipotent
// shutdown function.
func New(capacity int) (log.Logger, func(context.Context), chanlog.RecordsChannel) {
	exp, _ := chanlog.New(chanlog.WithCap(capacity))
	proc := sdklog.NewSimpleProcessor(exp)
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(proc))
	l := lp.Logger("testlogger",
		log.WithInstrumentationAttributes(
			attribute.Int(InstrumentationAttributeName, InstrumentationAttributeValue)))

	return l, func(ctx context.Context) { _ = lp.Shutdown(ctx) }, exp.Ch()
}
