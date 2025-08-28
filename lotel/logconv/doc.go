/*
Package logconv provides (the missing) conversions forth and back between
“plain” any values and OpenTelemetry's log values, including log key-values.

These conversions are intended to be used for testing, they are not meant for
production use. When writing tests, the conversions greatly reduce the otherwise
onerous OpenTelemetry [log.Value] and [log.KeyValue] setup boilerplate.

For example:

	v := logconv.Value([]any{"foo", 42, "bar"})
	a := logconv.Any(v)
	Expect(logconv.Value(a).Equal(v)).To(BeTrue())
*/
package logconv
