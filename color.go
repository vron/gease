package gease

import (
	"image/color"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

// ColorEasing smoothly animates a color transition. It operates in LAB space with
// seperate alpha interpolation for a visually smooth transition.
type ColorEasing struct {
	target      [3]float64
	targetAlpha float64

	l, a, b, alpha *interpolator

	last time.Time
}

// Color creates a new ColorEasing.
func Color(color color.RGBA) *ColorEasing {
	es := &ColorEasing{
		last: time.Now(),
	}
	es.targ(color)
	es.l = interp(es.target[0], 0.0, 1)
	es.a = interp(es.target[1], 0.0, 1)
	es.b = interp(es.target[2], 0.0, 1)
	es.alpha = interp(es.targetAlpha, 0.0, 1)
	es.Configure(DefaultPeriod)

	return es
}

// Configure sets the natural period of the animation. The provided
// period will thus roughly represent the time of a transition.
func (es *ColorEasing) Configure(period time.Duration) *ColorEasing {
	p := float64(period) / float64(time.Second)
	es.l.configure(0, p)
	es.a.configure(0, p)
	es.b.configure(0, p)
	es.alpha.configure(0, p)
	return es
}

// Target sets the target color we should animate towards.
func (es *ColorEasing) Target(target color.RGBA) {
	es.targ(target)
	if es.converged() {
		es.last = time.Now()
	}
}

func (es *ColorEasing) targ(target color.RGBA) {
	c := colorful.Color{
		R: float64(target.R) / float64(target.A),
		G: float64(target.G) / float64(target.A),
		B: float64(target.B) / float64(target.A),
	}
	es.target[0], es.target[1], es.target[2] = c.Lab()
	es.targetAlpha = float64(target.A) / float64(255)
}

// Step advances the simulation until t. If t is before any previos t Step was called with
// no work is done. If t.IsZero() time.Now() is assumed.
func (es *ColorEasing) Step(t time.Time) (converged bool) {
	if t.IsZero() {
		t = time.Now()
	}

	dt := float64(t.Sub(es.last)) / 1e9
	es.last = t

	es.l.step(es.target[0], dt)
	es.a.step(es.target[1], dt)
	es.b.step(es.target[2], dt)
	es.alpha.step(es.targetAlpha, dt)

	return es.converged()
}

func (es *ColorEasing) converged() bool {
	p := 0.25 / 255 / 255
	return es.l.converged(es.target[0], p) &&
		es.a.converged(es.target[1], p) &&
		es.b.converged(es.target[2], p) &&
		es.alpha.converged(es.targetAlpha, p)
}

// V returns the current color of the easing.
func (es *ColorEasing) V() color.RGBA {
	l := es.l.x
	a := es.a.x
	b := es.b.x
	alpha := es.alpha.x
	cr, cg, cb, _ := colorful.Lab(l, a, b).RGBA()
	alpha8 := uint16(alpha*65535 + 0.5)

	r := uint32(cr)
	r *= uint32(alpha8)
	r /= 0xffff
	r /= 257
	g := uint32(cg)
	g *= uint32(alpha8)
	g /= 0xffff
	g /= 257
	bb := uint32(cb)
	bb *= uint32(alpha8)
	bb /= 0xffff
	bb /= 257

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(bb),
		A: uint8(alpha*255 + 0.5),
	}
}

// SetTime stores the time t as the current time of the animation, without
// advancing the simulation. Usually there is no need to call SetTime since
// Color() internally set the initial time to time.Now().
func (es *ColorEasing) SetTime(t time.Time) {
	es.last = t
}
