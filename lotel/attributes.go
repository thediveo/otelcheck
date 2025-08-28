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
	sdklog "go.opentelemetry.io/otel/sdk/log"

	ty "github.com/onsi/gomega/types"
)

// attributeMatcher marks a Gomega matcher to match OTel log record attributes.
type attributeMatcher interface {
	matchAttribute(log.KeyValue) (bool, error)
}

// matchAllAttributes succeeds if the attributes of the passed log record match
// all passed attribute matchers.
func matchAllAttributes(r *sdklog.Record, ms []attributeMatcher) (bool, error) {
	for attr := range r.WalkAttributes /* sweet iterator */ {
		for _, m := range ms {
			success, err := m.matchAttribute(attr)
			if err != nil || !success {
				return false, err
			}
		}
	}
	return true, nil
}

// separateAttributeMatchers separates a list of matchers into a list of
// attribute matchers as well as the list of non-attribute matchers.
func separateAttributeMatchers(ms []ty.GomegaMatcher) ([]ty.GomegaMatcher, []attributeMatcher) {
	gms := make([]ty.GomegaMatcher, 0, len(ms))
	var ams []attributeMatcher
	for _, m := range ms {
		if am, ok := m.(attributeMatcher); ok {
			ams = append(ams, am)
			continue
		}
		gms = append(gms, m)
	}
	return gms, ams
}

func wrapAttributeMatcher(m ty.GomegaMatcher) ty.GomegaMatcher {
	return &wrappedMatcher{m}
}

type wrappedMatcher struct {
	ty.GomegaMatcher
}

var _ (attributeMatcher) = (*wrappedMatcher)(nil)

func (m *wrappedMatcher) matchAttribute(attr log.KeyValue) (bool, error) {
	return m.Match(attr)
}
