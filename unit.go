package gease

import (
	"time"

	"gioui.org/unit"
)

// UnitEasing soothly animates a gioui unit.Value. Note that the unit (e..g px, sp, dp)
// must be consistent in each UnitEasing.
type UnitEasing struct {
	target unit.Value

	x *interpolator

	last time.Time
}

// Unit creaes a new UnitEasing.
func Unit(v unit.Value) *UnitEasing {
	es := &UnitEasing{target: v,
		x:    interp(float64(v.V), 0.0, 1.0),
		last: time.Now(),
	}
	es.Configure(DefaultOvershoot, DefaultPeriod)
	return es
}

// Configure sets the stiffness of the animation by specifing overshoot and period. The provided
// period will thus roughly represent the time of a transition abd overshoot the amount
// by which the animation should go past the target value before returning.
func (es *UnitEasing) Configure(overshoot float64, period time.Duration) *UnitEasing {
	p := float64(period) / float64(time.Second)
	es.x.configure(overshoot, p)
	return es
}

// Target sets the target unit.Value we should animate to.
func (es *UnitEasing) Target(target unit.Value) {
	if target.U != es.target.U {
		// panic("UnitEasing must be used with consistent units")

		// instead of panicing and possibly killing a client application we will simply
		// move on, assuming the original unit. This might lead to unexpected values but
		// is clearly documented that the user should not do.
		target.U = es.target.U
	}
	es.target = target

	if es.converged() {
		es.last = time.Now()
	}
}

// Step advances the simulation until t. If t is before any previos t Step was called with
// the current value is returned without advancing the animation. If t.IsZero() time.Now() is assumed.
func (es *UnitEasing) Step(t time.Time) (converged bool) {
	if t.IsZero() {
		t = time.Now()
	}

	dt := float64(t.Sub(es.last)) / 1e9
	es.last = t
	es.x.step(float64(es.target.V), dt)

	return es.converged()
}

func (es *UnitEasing) converged() bool {
	p := 0.1
	return es.x.converged(float64(es.target.V), p)
}

// V returns the current value of the easing.
func (es *UnitEasing) V() unit.Value {
	return unit.Value{U: es.target.U, V: float32(es.x.x)}
}

// SetTime stores the time t as the current time of the animation, without
// advancing the simulation. Usually there is no need to call SetTime since
// Color() internally set the initial time to time.Now().
func (es *UnitEasing) SetTime(t time.Time) {
	es.last = t
}
