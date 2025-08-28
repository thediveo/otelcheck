/*
Package chanlog provides an exporter for OpenTelemetry log telemetry that sends
log records into a Go channel.

This exporter is intended to be used for testing, it is not meant for production
use.

Tests then can pick up the log records emitted by the code under test from the
Go channel, either concurrently or at certain check points, leveraging channel
buffering.
*/
package chanlog
