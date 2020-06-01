package gease

import (
	"math"
	"testing"
	"time"

	"gioui.org/unit"
)

func TestConverge(t *testing.T) {
	hub := NewHub()
	es1 := Unit(unit.Dp(10))
	es2 := Unit(unit.Dp(20))
	es3 := Unit(unit.Dp(30))

	es1.Target(unit.Dp(33))
	es2.Target(unit.Dp(33))
	es3.Target(unit.Dp(33))

	hub.Add(es1, es2, es3)
	hub.Remove(es1)

	tt := time.Now().Add(time.Millisecond * 10)
	for !hub.Step(tt) {
		tt = tt.Add(time.Millisecond * 10)
	}

	// We expect it to have converged on both these two, but also
	// that the first one ahs not been updated.
	if es1.V() != unit.Dp(10) {
		t.Error("expected not changed")
	}

	if math.Abs(float64(es2.V().V-33)) > 0.15 || math.Abs(float64(es3.V().V-33)) > 0.15 {
		t.Error("expected them to change")
	}
}
