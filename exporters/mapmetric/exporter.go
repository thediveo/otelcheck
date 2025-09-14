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

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type Exporter struct {
	temporalitySelector metric.TemporalitySelector
}

var _ metric.Exporter = (*Exporter)(nil)

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
	e := &Exporter{}
	return e, nil
}

func (e *Exporter) Aggregation(metric.InstrumentKind) metric.Aggregation {}

func (e *Exporter) Temporality(metric.InstrumentKind) metricdata.Temporality {}

func (e *Exporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error {
}

func (e *Exporter) ForceFlush(context.Context) error { return nil }

func (e *Exporter) Shutdown(ctx context.Context) error {
}
