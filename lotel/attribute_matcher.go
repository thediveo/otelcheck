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
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"

	"github.com/thediveo/otelcheck/lotel/logconv"

	"github.com/onsi/gomega/format"
	ty "github.com/onsi/gomega/types"
)

// HaveAttribute succeeds if an OpenTelemetry log record has the an attribute
// with the specified key/name (and optional value).
//
// The expected value passed into the attr parameter can be either a string or
// [ty.GomegaMatcher]:
//   - a string in the form of “name” where it must match an attribute key/name
//     (but not any value), or
//     in the “name=value” form where it must match both the attribute key/name and
//     value. Please note that the “name=value” form matches attribute string values
//     only. Use [HaveAttributeWIthValue] instead to match attributes with empty
//     or non-string values.
//   - a GomegaMatcher that matches the name only.
//   - any other type of value is an error.
//
// HaveAttribute accepts actual values of types [log.KeyValue] and also
// [sdklog.Record]. When actual is a log record as opposed to a key-value pair,
// HaveAttribute matches against the attribute set of this log record. However,
// it is recommended to use HaveAttribute only in the context of BeARecord,
// especially when asserting the presence of multiple attributes, as BeARecord
// optimizes its attributes checks.
//
// Usage examples:
//
//	HaveAttribute("foo")
//	HaveAttribute("foo=bar")
//	HaveAttribute(HaveSuffix("foo"))
//
// See also [HaveAttributeWithValue].
func HaveAttribute(attr any) ty.GomegaMatcher {
	if s, ok := attr.(string); ok {
		// single plain string argument, so let's see if it is in "NAME=VALUE"
		// format...
		if nam, value, found := strings.Cut(s, "="); found {
			return &HaveAttributeMatcher{
				name:         nam,
				value:        value,
				nameMatcher:  matcherOrEqual(nam),
				valueMatcher: matcherOrEqual(value),
			}
		}
		// it's just "NAME", so we want to match only the name, not the value;
		// fall through now.
	}
	return &HaveAttributeMatcher{
		name:         attr,
		value:        nil,
		nameMatcher:  matcherOrEqual(attr),
		valueMatcher: nil,
	}
}

// HaveAttributeWithValue succeeds if an OpenTelemetry log record has an
// attribute with the specified key/name and value.
//
// The value passed into the name parameter can be either a string or a
// [ty.GomegaMatcher].
//
// HaveAttributeWithValue accepts actual values of types [log.KeyValue] and also
// [sdklog.Record]. When actual is a log record as opposed to a key-value pair,
// HaveAttributeWithValue matches against the attribute set of this log record.
// However, it is recommended to use HaveAttributeWithValue only in the context
// of BeARecord, especially when asserting the presence of multiple attributes,
// as BeARecord optimizes the attributes checks.
//
// The value passed into the value parameter can be one of the following, all
// other values are an error:
//   - nil: represents an OTel “empty” value.
//   - bool
//   - int, int64
//   - float32, float64
//   - string
//   - []byte
//   - []any: represents OTel slice values
//   - map[string]any: represents OTel map values
//   - [ty.GomegaMatcher]
//
// Usage examples:
//
//	HaveAttributeWithValue("foo", "bar")
//	HaveAttributeWithValue("foo", nil) // explicitly check for empty-ness
//	HaveAttributeWithValue("foo", Not(BeEmpty()))
//	HaveAttributeWithValue("foo", 42)
//
// See also [HaveAttribute].
func HaveAttributeWithValue(name, value any) ty.GomegaMatcher {
	return &HaveAttributeMatcher{
		name:         name,
		value:        value,
		nameMatcher:  matcherOrEqual(name),
		valueMatcher: matcherOrEqualNilInclusive(value, logconv.Canonize),
	}

}

type HaveAttributeMatcher struct {
	name         any
	value        any
	nameMatcher  ty.GomegaMatcher
	valueMatcher ty.GomegaMatcher
}

var (
	_ attributeMatcher = (*HaveAttributeMatcher)(nil)
	_ ty.GomegaMatcher = (*HaveAttributeMatcher)(nil)
)

func (m *HaveAttributeMatcher) matchAttribute(attr log.KeyValue) (bool, error) {
	if m.nameMatcher == nil {
		return false, fmt.Errorf("HaveAttributeMatcher: name matcher must not be <nil>")
	}
	if m.value != nil && m.valueMatcher == nil {
		return false, fmt.Errorf("HaveAttributeMatcher: expected value to match to be either string or types.GomegaMatcher.  Got:\n%T",
			m.value)
	}
	success, err := m.nameMatcher.Match(attr.Key)
	if err != nil || !success {
		return false, err
	}
	if m.valueMatcher == nil { // no value to match, so we've found a matching attribute
		return true, nil
	}
	return m.valueMatcher.Match(logconv.Any(attr.Value))
}

func (m *HaveAttributeMatcher) Match(actual any) (success bool, err error) {
	if actual == nil {
		return false, errors.New("refusing to match <nil>")
	}
	r, ok := actual.(sdklog.Record)
	if ok {
		for attr := range r.WalkAttributes {
			if success, err := m.matchAttribute(attr); err != nil || success {
				return success, err
			}
		}
		return false, nil
	}
	attr, ok := actual.(log.KeyValue)
	if !ok {
		return false, fmt.Errorf("HaveAttribute expected actual of type <%T> or <%T>.  Got:\n%s",
			log.KeyValue{}, sdklog.Record{}, format.Object(actual, 1))
	}
	return m.matchAttribute(attr)
}

func (m *HaveAttributeMatcher) expected() string {
	expected := "key:\n" + format.Object(m.name, 1)
	if m.value != nil {
		expected += "\nvalue:\n" + format.Object(m.value, 1)
	}
	return expected
}

func (m *HaveAttributeMatcher) FailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n%s\nto equal\n%s",
		format.Object(actual, 1), format.IndentString(m.expected(), 1))
}

func (m *HaveAttributeMatcher) NegatedFailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n%s\nnot to equal\n%s",
		format.Object(actual, 1), format.IndentString(m.expected(), 1))
}
