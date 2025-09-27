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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var _ = Describe("metrics indices", func() {

	It("rejects the same-name metrics with differing aggregations", func() {
		rms := []*metricdata.ResourceMetrics{
			{
				ScopeMetrics: []metricdata.ScopeMetrics{
					{
						Metrics: []metricdata.Metrics{
							{
								Name: "foo.bar",
								Data: metricdata.Gauge[int64]{
									DataPoints: []metricdata.DataPoint[int64]{
										{Value: 42},
									},
								},
							},
							{
								Name: "foo.bar",
								Data: metricdata.Sum[int64]{
									DataPoints: []metricdata.DataPoint[int64]{
										{Value: 666},
									},
								},
							},
						},
					},
				},
			},
		}
		Expect(Index(rms)).Error().To(MatchError(MatchRegexp(
			`metric "foo.bar" .* Sum\[int64\] and Gauge\[int64\]`)))
	})

})
