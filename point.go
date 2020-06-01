package gease

import (
	"time"

	"gioui.org/f32"
)

// PointEasing soothly animates a 2D-position.
type PointEasing struct {
	target f32.Point

	x, y *interpolator

	last time.Time
}

// Point creaes a new PointEasing.
func Point(pos f32.Point) *PointEasing {
	es := &PointEasing{target: pos,
		x:    interp(float64(pos.X), 0.0, 1.0),
		y:    interp(float64(pos.Y), 0.0, 1.0),
		last: time.Now(),
	}
	es.Configure(DefaultOvershoot, DefaultPeriod)
	return es
}

// Configure sets the stiffness of the animation by specifing overshoot and period. The provided
// period will thus roughly represent the time of a transition abd overshoot the amount
// by which the animation should go past the target value before returning.
func (es *PointEasing) Configure(overshoot float64, period time.Duration) *PointEasing {
	p := float64(period) / float64(time.Second)
	es.x.configure(overshoot, p)
	es.y.configure(overshoot, p)
	return es
}

// Target sets the target point we should animate to.
func (es *PointEasing) Target(target f32.Point) {
	es.target = target

	if es.converged() {
		es.last = time.Now()
	}
}

// Step advances the simulation until t. If t is before any previos t Step was called with
// no work is done. If t.IsZero() time.Now() is assumed.
func (es *PointEasing) Step(t time.Time) (converged bool) {
	if t.IsZero() {
		t = time.Now()
	}

	dt := float64(t.Sub(es.last)) / 1e9
	es.last = t

	es.x.step(float64(es.target.X), dt)
	es.y.step(float64(es.target.Y), dt)

	return es.converged()
}

func (es *PointEasing) converged() bool {
	p := 0.1
	return es.x.converged(float64(es.target.X), p) &&
		es.y.converged(float64(es.target.Y), p)
}

// V returns the current point of the easing.
func (es *PointEasing) V() f32.Point {
	return f32.Point{X: float32(es.x.x),
		Y: float32(es.y.x)}
}

// SetTime stores the time t as the current time of the animation, without
// advancing the simulation. Usually there is no need to call SetTime since
// Color() internally set the initial time to time.Now().
func (es *PointEasing) SetTime(t time.Time) {
	es.last = t
}
