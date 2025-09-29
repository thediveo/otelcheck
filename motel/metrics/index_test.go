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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var _ = Describe("metrics indices", func() {

	Context("same-name metrics with differing properties", func() {

		It("rejects differing units", func() {
			rms := []*metricdata.ResourceMetrics{
				{
					ScopeMetrics: []metricdata.ScopeMetrics{
						{
							Metrics: []metricdata.Metrics{
								{
									Name: "foo.bar",
									Unit: "kilofubars",
									Data: metricdata.Gauge[int64]{
										DataPoints: []metricdata.DataPoint[int64]{
											{
												Value:      42,
												Attributes: attribute.NewSet(attribute.String("foo", "bar")),
											},
										},
									},
								},
								{
									Name: "foo.bar",
									Unit: "megafubars",
									Data: metricdata.Sum[int64]{
										DataPoints: []metricdata.DataPoint[int64]{
											{
												Value:      42,
												Attributes: attribute.NewSet(attribute.String("foo", "baz")),
											},
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(Index(rms)).Error().To(MatchError(MatchRegexp(
				`metric "foo.bar" with differing units "megafubars" and "kilofubars"`)))
		})

		It("rejects differing descriptions", func() {
			rms := []*metricdata.ResourceMetrics{
				{
					ScopeMetrics: []metricdata.ScopeMetrics{
						{
							Metrics: []metricdata.Metrics{
								{
									Name:        "foo.bar",
									Description: "The lazy brown fox",
									Data: metricdata.Gauge[int64]{
										DataPoints: []metricdata.DataPoint[int64]{
											{
												Value:      42,
												Attributes: attribute.NewSet(attribute.String("foo", "bar")),
											},
										},
									},
								},
								{
									Name:        "foo.bar",
									Description: "The fury lazy socks",
									Data: metricdata.Sum[int64]{
										DataPoints: []metricdata.DataPoint[int64]{
											{
												Value:      42,
												Attributes: attribute.NewSet(attribute.String("foo", "baz")),
											},
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(Index(rms)).Error().To(MatchError(MatchRegexp(
				`metric "foo.bar" with differing descriptions ".+" and ".+"`)))
		})

		It("rejects differing aggregations", func() {
			rms := []*metricdata.ResourceMetrics{
				{
					ScopeMetrics: []metricdata.ScopeMetrics{
						{
							Metrics: []metricdata.Metrics{
								{
									Name: "foo.bar",
									Data: metricdata.Gauge[int64]{
										DataPoints: []metricdata.DataPoint[int64]{
											{
												Value:      42,
												Attributes: attribute.NewSet(attribute.String("foo", "bar")),
											},
										},
									},
								},
								{
									Name: "foo.bar",
									Data: metricdata.Sum[int64]{
										DataPoints: []metricdata.DataPoint[int64]{
											{
												Value:      42,
												Attributes: attribute.NewSet(attribute.String("foo", "baz")),
											},
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

	It("collects and indexes metrics into families", func() {
		rms := []*metricdata.ResourceMetrics{
			{
				ScopeMetrics: []metricdata.ScopeMetrics{
					{
						Metrics: []metricdata.Metrics{
							{
								Name: "foo.bar",
								Data: metricdata.Gauge[int64]{
									DataPoints: []metricdata.DataPoint[int64]{
										{
											Value:      42,
											Attributes: attribute.NewSet(attribute.String("foo", "bar")),
										},
									},
								},
							},
						},
					},
					{
						Metrics: []metricdata.Metrics{
							{
								Name: "foo.bar",
								Data: metricdata.Gauge[int64]{
									DataPoints: []metricdata.DataPoint[int64]{
										{
											Value:      666,
											Attributes: attribute.NewSet(attribute.String("foo", "baz")),
										},
									},
								},
							},
						},
					},
				},
			},
			{
				ScopeMetrics: []metricdata.ScopeMetrics{
					{
						Metrics: []metricdata.Metrics{
							{
								Name: "foo.bar",
								Data: metricdata.Gauge[int64]{
									DataPoints: []metricdata.DataPoint[int64]{
										{
											Value:      123,
											Attributes: attribute.NewSet(attribute.String("foo", "tada")),
										},
									},
								},
							},
						},
					},
				},
			},
		}
		Expect(Index(rms)).To(HaveKeyWithValue("foo.bar",
			ConsistOf(
				HaveField("Metrics.Data.DataPoints", ConsistOf(
					HaveField("Value", int64(42)))),
				HaveField("Metrics.Data.DataPoints", ConsistOf(
					HaveField("Value", int64(123)))),
				HaveField("Metrics.Data.DataPoints", ConsistOf(
					HaveField("Value", int64(666)))))))
	})
})
