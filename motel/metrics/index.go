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

package metrics

import (
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type MetricsIndex map[string][]Metrics

type Metrics struct {
	*metricdata.Metrics
	Resource *resource.Resource
	Scope    *instrumentation.Scope
}

func Index(rms /* ...no-no, not that one, we're Apache here */ []*metricdata.ResourceMetrics) (MetricsIndex, error) {
	index := MetricsIndex{}
	for _, rm := range rms {
		for sidx, sm := range rm.ScopeMetrics {
			for midx, m := range sm.Metrics {
				ms, ok := index[m.Name]
				if !ok {
					// Never seen any metric of this name, so this is a first:
					// just index it and we're good.
					index[m.Name] = []Metrics{
						{
							Metrics:  &sm.Metrics[midx],
							Resource: rm.Resource,
							Scope:    &rm.ScopeMetrics[sidx].Scope,
						},
					}
					continue
				}
				// We've seen other metrics with this name, so check that we're
				// consistent with what we've already seen...
				if m.Unit != ms[0].Unit {
					return nil, fmt.Errorf("metric %q with differing units %q and %q",
						m.Name, m.Unit, ms[0].Unit)
				}
				if m.Description != ms[0].Description {
					return nil, fmt.Errorf("metric %q with differing descriptions %q and %q",
						m.Name, m.Description, ms[0].Description)
				}
				d1t := reflect.TypeOf(m.Data)
				d2t := reflect.TypeOf(ms[0].Data)
				if d1t != d2t {
					return nil, fmt.Errorf("metric %q with differing aggregations %s and %s",
						m.Name, d1t.Name(), d2t.Name())
				}

				index[m.Name] = append(ms, Metrics{
					Metrics:  &sm.Metrics[midx],
					Resource: rm.Resource,
					Scope:    &rm.ScopeMetrics[sidx].Scope,
				})
			}
		}
	}
	return index, nil
}
