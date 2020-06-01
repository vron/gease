package gease

import (
	"testing"
	"time"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

func TestColorEnds(t *testing.T) {
	red, green := colornames.Red400, colornames.Green400
	es := Color(red)
	es.Target(green)

	// check initial color
	if es.V() != red {
		t.Log(red, es.Step(time.Time{}))
		t.Error("expected initial color at time 0")
	}

	// advance the simulation
	tt := time.Now()
	for !es.Step(tt) {
		tt = tt.Add(time.Millisecond * 10)
	}

	// check the final color
	if es.V() != green {
		t.Log(green, es.Step(time.Time{}))
		t.Error("expected final color at end")
	}
}

func TestColorLargestep(t *testing.T) {
	es := Color(colornames.Red400)
	es.Target(colornames.White)

	// check the final color
	es.Step(time.Now().Add(time.Hour))
	col := es.V()
	if col != colornames.White {
		t.Error("largestep not working")
	}
}

func BenchmarkColorStep(b *testing.B) {
	es := Color(colornames.White)
	es.Target(colornames.Black)
	b.ReportAllocs()
	t := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		es.Step(t)
		t = t.Add(time.Millisecond)
	}

	// prevent compiler removing
	b.StopTimer()
	if es.V().A == 128 {
		panic("cannot happen")
	}
}

func BenchmarkColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		es := Color(colornames.White)
		es.Target(colornames.Black)

		// assuming 60 fps here
		t := time.Now()
		for !es.Step(t) {
			t = t.Add(time.Millisecond * 17)
		}
		// prevent compiler removing
		if es.V().A == 128 {
			panic("cannot happen")
		}
	}
}
