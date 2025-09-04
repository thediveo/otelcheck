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

// HaveSeverityText succeeds if the actual log record has the expected severity
// text. The expected severity text can either be a string or alternatively a
// [ty.GomegaMatcher].
func HaveSeverityText(expected any) ty.GomegaMatcher {
	m := matcherOrEqual(expected)
	return gc.MakeMatcher(func(r sdklog.Record) (bool, error) {
		return m.Match(r.SeverityText())
	}).WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} match\n{{format .Data 1}}").
		WithTemplateData(expected)
}
