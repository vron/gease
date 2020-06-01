package gease

import (
	"math"
	"testing"
	"time"

	"gioui.org/unit"
)

func TestUnitEnds(t *testing.T) {
	target := unit.Dp(10)
	es := Unit(unit.Dp(4))
	es.Target(target)

	// check initial position
	if es.V() != unit.Dp(4) {
		t.Log(es.Step(time.Time{}))
		t.Error("expected initial")
		t.FailNow()
	}

	// advance the simulation
	tt := time.Now()
	for !es.Step(tt) {
		tt = tt.Add(time.Millisecond * 270)
	}

	// check the final position
	pos := es.V()
	if math.Abs(float64(target.V-pos.V)) > 0.15 {
		t.Error("expected final position at end")
	}
}

func BenchmarkUnitStep(b *testing.B) {
	es := Unit(unit.Dp(10))
	es.Target(unit.Dp(4))
	b.ReportAllocs()
	t := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		es.Step(t)
		t = t.Add(time.Millisecond)
	}

	// prevent compiler removing
	b.StopTimer()
	if es.V().V == -33 {
		panic("cannot happen")
	}
}

func BenchmarkUnit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		es := Unit(unit.Dp(4))
		es.Target(unit.Dp(10))

		// assuming 60 fps here
		t := time.Now()
		for !es.Step(t) {
			t = t.Add(time.Millisecond * 17)
		}
		// prevent compiler removing
		if es.V().V == 128 {
			panic("cannot happen")
		}
	}
}
