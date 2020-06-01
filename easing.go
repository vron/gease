/*
Package gease implements spring based animation easings for gioui

Spring / force based animation of GUI components for gioui. In general
plenty fast enough to always drive all dimensions through this package for
smooth and nice animations at e.g. resizes etc.

Please see github.com/vron/gease/example as an example of the intended use-case.

*/
package gease

import (
	"time"
)

const (
	// DefaultOvershoot represents a resonable default for visualizations.
	DefaultOvershoot = 0.1
	// DefaultPeriod represents a reasonable default for visualizations.
	DefaultPeriod = 750 * time.Millisecond
)
