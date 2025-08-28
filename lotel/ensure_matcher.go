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
	"go.opentelemetry.io/otel/log"

	g "github.com/onsi/gomega"
	ty "github.com/onsi/gomega/types"
	"github.com/thediveo/otelcheck/lotel/logconv"
)

// matcherOrEqual either returns any passed-in [ty.GomegaMatcher] value as-is
// and otherwise wraps all other expected values into a [g.Equal] matcher.
//
// See also: [toMatcherNilInclusive]
func matcherOrEqual(expected any) ty.GomegaMatcher {
	if m, ok := expected.(ty.GomegaMatcher); ok {
		return m
	}
	return g.Equal(expected)
}

// matcherOrEqualNilInclusive either returns any passed-in [ty.GomegaMatcher] value
// as-is and otherwise wraps all other expected values into either a [g.Equal]
// or [g.BeNil] matcher, depending on expected. The dedicated handling of
// expected nil values allows to match “empty” log values (which we represent as
// nil after any-fying [log.Value] to any values).
func matcherOrEqualNilInclusive(expected any, fn ...func(any) any) ty.GomegaMatcher {
	if m, ok := expected.(ty.GomegaMatcher); ok {
		return m
	}
	if expected == nil {
		return g.BeNil()
	}
	if len(fn) > 0 {
		return g.Equal(fn[0](expected))
	}
	return g.Equal(expected)
}

// valueMatcher either returns any passed-in [ty.GomegaMatcher] value as is
// and otherwise wraps all other expected values into a [EqualsValue] matcher,
// translating the expected(!) value into a log value, where necessary.
func valueMatcher(expected any) ty.GomegaMatcher {
	if m, ok := expected.(ty.GomegaMatcher); ok {
		return m
	}
	if expected, ok := expected.(log.Value); ok {
		return EqualsValue(expected)
	}
	return EqualsValue(logconv.Value(expected))
}
