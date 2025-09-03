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
	sdklog "go.opentelemetry.io/otel/sdk/log"

	gc "github.com/onsi/gomega/gcustom"
	ty "github.com/onsi/gomega/types"
)

// BeARecord succeeds if the actual log record satisfies all specified matchers.
// It is an error for actual not to be of type [sdklog.Record].
func BeARecord(m ty.GomegaMatcher, ms ...ty.GomegaMatcher) ty.GomegaMatcher {
	ms = append([]ty.GomegaMatcher{m}, ms...)
	gms, ams := separateAttributeMatchers(ms)
	return gc.MakeMatcher(func(r sdklog.Record) (bool, error) {
		for _, m := range gms {
			success, err := m.Match(r)
			if err != nil || !success {
				return false, err
			}
		}
		return containsAttributes(&r, ams)
	}).WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} match\n{{format .Data 1}}").
		WithTemplateData(ms)
}
