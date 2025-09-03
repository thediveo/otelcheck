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
	"slices"

	"go.opentelemetry.io/otel/attribute"
	sdklog "go.opentelemetry.io/otel/sdk/log"

	intslices "github.com/thediveo/otelcheck/internal/slices"
	"github.com/thediveo/otelcheck/lotel/logconv"

	ty "github.com/onsi/gomega/types"
)

// attributeMatcher marks a Gomega matcher to match OTel log record attributes.
type attributeMatcher interface {
	// try to match an attribute by its name and any-fied/canonized value, where
	// the attribute name/value might come from a log.KeyValue or resource/scope
	// attribute.KeyValue.
	matchAttribute(name string, value any) (bool, error)
}

// containsAttributes succeeds if all passed attribute matchers match on (some
// of) the passed log record's attributes including resource and scope
// attributes.
func containsAttributes(r *sdklog.Record, attrms []attributeMatcher) (bool, error) {
	attrms, err := removeMatchingMatchers(r.Resource().Set(), slices.Clone(attrms))
	if err != nil {
		return false, err
	}
	if len(attrms) == 0 {
		return true, nil
	}
	attrs := r.InstrumentationScope().Attributes
	attrms, err = removeMatchingMatchers(&attrs, attrms)
	if err != nil {
		return false, err
	}
	if len(attrms) == 0 {
		return true, nil
	}
	// And now, esteemed brethren, we enter the last chance saloon...
nextRecordAttribute:
	for attr := range r.WalkAttributes /* sweet iterator */ {
		if len(attrms) == 0 {
			return true, nil
		}
		key := attr.Key
		value := logconv.Any(attr.Value)
		for midx, m := range attrms {
			success, err := m.matchAttribute(key, value)
			if err != nil {
				return false, err
			}
			if success {
				attrms = intslices.DeleteUnordered(attrms, midx)
				continue nextRecordAttribute
			}
		}
	}
	return len(attrms) == 0, nil
}

// removeMatchingMatchers checks which attribute matchers match on the
// passed attribute set and then returns only the "left-over" non-matching
// matchers.
func removeMatchingMatchers(attrs *attribute.Set, attrms []attributeMatcher) ([]attributeMatcher, error) {
	it := attrs.Iter()
nextAttribute:
	for it.Next() {
		if len(attrms) == 0 {
			return attrms, nil
		}
		attr := it.Attribute()
		value := logconv.Canonize(attr.Value.AsInterface())
		for midx, m := range attrms {
			success, err := m.matchAttribute(string(attr.Key), value)
			if err != nil {
				return nil, err
			}
			if success {
				attrms = intslices.DeleteUnordered(attrms, midx)
				continue nextAttribute
			}
		}
	}
	return attrms, nil
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
