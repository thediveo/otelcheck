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

package mapmetric

import (
	"context"
	"sync/atomic"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type Exporter struct {
	temporalitySelector metric.TemporalitySelector
	aggregationSelector metric.AggregationSelector

	mets atomic.Pointer[[]*metricdata.ResourceMetrics]
}

// statically ensure that we fulfill the OTel metrics SDK's Exporter interface.
var _ metric.Exporter = (*Exporter)(nil)

// New returns a new metrics exporter, configured with the passed options.
func New(opts ...Option) (*Exporter, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	if o.temporalitySelector == nil {
		o.temporalitySelector = metric.DefaultTemporalitySelector
	}
	if o.aggregationSelector == nil {
		o.aggregationSelector = metric.DefaultAggregationSelector
	}
	e := &Exporter{
		temporalitySelector: o.temporalitySelector,
		aggregationSelector: o.aggregationSelector,
	}
	return e, nil
}

// Aggregation returns the Aggregation to use for an instrument kind.
func (e *Exporter) Aggregation(k metric.InstrumentKind) metric.Aggregation {
	return e.aggregationSelector(k)
}

// Temporality returns the Temporality to use for an instrument kind.
func (e *Exporter) Temporality(k metric.InstrumentKind) metricdata.Temporality {
	return e.temporalitySelector(k)
}

// Export metrics to our internal slice storage. It does nothing after
// [Exporter.Shutdown] has been called.
func (e *Exporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	mets := e.mets.Load()
	if mets == nil {
		return ctx.Err()
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	*mets = append(*mets, rm)
	return nil
}

// Shutdown the Exporter so that any later calls to [Exporter.Export] will
// perform no operation anymore.
func (e *Exporter) Shutdown(ctx context.Context) error {
}

// ForceFlush is a no-op.
func (e *Exporter) ForceFlush(context.Context) error { return nil }
