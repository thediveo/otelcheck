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

package slices

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("slices", func() {

	DescribeTable("delete unordered",
		func(s []string, idx int, expected []string) {
			Expect(DeleteUnordered(s, idx)).To(ConsistOf(expected))
		},
		Entry(nil, []string{"foo", "bar", "baz"}, 0, []string{"bar", "baz"}),
		Entry(nil, []string{"foo", "bar", "baz"}, 1, []string{"foo", "baz"}),
		Entry(nil, []string{"foo", "bar", "baz"}, 2, []string{"foo", "bar"}),
		Entry(nil, []string{"foo"}, 0, []string{}),
	)

})
