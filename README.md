# gease
Easings / animations for gioui

![Test](https://github.com/vron/gease/workflows/Test/badge.svg)

WIP: looking for input on API / feedback.

This package implements spring based convenience methods for animating layouts or draw operations
in Gioui. Please refer to https://godoc.org/github.com/vron/gease for further details.

To see smooth easings on positions, sizes and colors you can run:

    go run github.com/vron/gease/example

Inspired by: https://www.react-spring.io/

## Allocations
The pacakge has beed designed to minimize GC pressure during animations, in particular there
is no additoinal allocation per frame resulting from using this easing package:

    go test -bench Step github.com/vron/gease       
    goos: windows
    goarch: amd64
    pkg: github.com/vron/gease
    BenchmarkColorStep-24           22201993                54.5 ns/op             0 B/op          0 allocs/op
    BenchmarkPointStep-24           44404644                26.0 ns/op             0 B/op          0 allocs/op
    BenchmarkUnitStep-24            59919408                19.6 ns/op             0 B/op          0 allocs/op
    PASS
    ok      github.com/vron/gease   3.799s