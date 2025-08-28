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

package logconv

import (
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel/log"
)

// Any returns the passed log value as an any value, recursively any-fying
// slice and map log values.
//
// Any returns slice values as []any and map values as map[string]any.
func Any(v log.Value) any {
	switch v.Kind() {
	case log.KindBool:
		return v.AsBool()
	case log.KindInt64:
		return v.AsInt64()
	case log.KindFloat64:
		return v.AsFloat64()
	case log.KindString:
		return v.AsString()
	case log.KindBytes:
		return v.AsBytes()
	case log.KindSlice:
		return anySlice(v)
	case log.KindMap:
		return anyMap(v)
	case log.KindEmpty:
		return nil
	}
	return nil
}

// anySlice returns the passed log slice value recursively any-fied.
func anySlice(v log.Value) []any {
	vs := v.AsSlice()
	sl := make([]any, 0, len(vs))
	for _, v := range vs {
		sl = append(sl, Any(v))
	}
	return sl
}

// anyMap returns the passed log map value recursively any-fied.
func anyMap(v log.Value) map[string]any {
	vm := v.AsMap()
	m := map[string]any{}
	for _, kv := range vm {
		m[kv.Key] = Any(kv.Value)
	}
	return m
}

// Canonize an any value to the allowed OpenTelemetry log.Value types by
// (recursively) promoting int to int64 and float32 to float64.
//
// Canonize panics when any value encountered is not of the following types:
//   - bool
//   - int and int64
//   - float32 and float64
//   - string
//   - []byte
//   - []any
//   - map[string]any
func Canonize(v any) any {
	if v == nil {
		return nil
	}
	switch v := v.(type) {
	case bool:
		return v
	case int:
		return int64(v)
	case int64:
		return v
	case float32:
		return float64(v)
	case float64:
		return v
	case string:
		return v
	case []byte:
		return v
	case []any:
		sl := make([]any, len(v))
		for idx, el := range v {
			sl[idx] = Canonize(el)
		}
		return sl
	case map[string]any:
		m := make(map[string]any, len(v))
		for key, value := range v {
			m[key] = Canonize(value)
		}
		return m
	}
	panic(fmt.Sprintf("logconv.Canonize: unsupported type %T", v))
}

// Value returns the log value for the passed (any) value.
//
// It panics for value types not supported by OTel's log value type.
func Value(v any) log.Value {
	if v == nil {
		return log.Value{} // KindEmpty
	}
	switch v := v.(type) {
	case bool:
		return log.BoolValue(v)
	case int:
		return log.IntValue(v)
	case int64:
		return log.Int64Value(v)
	case float32:
		return log.Float64Value(float64(v))
	case float64:
		return log.Float64Value(v)
	case string:
		return log.StringValue(v)
	case []byte:
		return log.BytesValue(v)
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		return valueSlice(v)
	case reflect.Map:
		return valueMap(v)
	}
	panic(fmt.Sprintf("logconv.Value: unsupported type %T", v))
}

func valueSlice(v any) log.Value {
	rv := reflect.ValueOf(v)
	l := rv.Len()
	vs := make([]log.Value, 0, l)
	for i := range l {
		vs = append(vs, Value(rv.Index(i).Interface()))
	}
	return log.SliceValue(vs...)
}

func valueMap(v any) log.Value {
	rv := reflect.ValueOf(v)
	kvs := make([]log.KeyValue, 0, rv.Len())
	mit := rv.MapRange()
	for mit.Next() {
		kvs = append(kvs, log.KeyValue{
			Key:   mit.Key().String(),
			Value: Value(mit.Value().Interface()),
		})
	}
	return log.MapValue(kvs...)
}
