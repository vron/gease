package gease

import (
	"math"
	"testing"
	"time"

	"gioui.org/f32"
)

func TestPointEnds(t *testing.T) {
	target := f32.Point{X: 100, Y: -55}
	es := Point(f32.Point{})
	es.Target(target)

	// check initial position
	if (es.V() != f32.Point{}) {
		t.Error("expected initial")
	}

	// advance the simulation
	tt := time.Now()
	for !es.Step(tt) {
		tt = tt.Add(time.Millisecond * 10)
	}

	// check the final position
	pos := es.V()
	if math.Abs(float64(pos.X-target.X)) > 0.15 || math.Abs(float64(pos.Y-target.Y)) > 0.15 {
		t.Error("expected final position at end")
	}
}

func BenchmarkPointStep(b *testing.B) {
	es := Point(f32.Point{})
	es.Target(f32.Point{X: 33, Y: 111})
	b.ReportAllocs()
	t := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		es.Step(t)
		t = t.Add(time.Millisecond)
	}

	// prevent compiler removing
	b.StopTimer()
	if es.V().X == -33 {
		panic("cannot happen")
	}
}

func BenchmarkPoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		es := Point(f32.Point{})
		es.Target(f32.Point{X: 33, Y: 111})

		// assuming 60 fps here
		t := time.Now()
		for !es.Step(t) {
			t = t.Add(time.Millisecond * 17)
		}
		// prevent compiler removing
		if es.V().X == 128 {
			panic("cannot happen")
		}
	}
}
