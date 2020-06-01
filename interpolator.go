package gease

import "math"

type interpolator struct {
	damp, omega float64
	x, v        float64
}

func interp(x float64, overshoot float64, period float64) *interpolator {
	i := &interpolator{
		x: x,
	}
	i.configure(overshoot, period)
	return i
}

func (i *interpolator) configure(overshoot float64, period float64) {
	if overshoot < 0 {
		overshoot = 0
	}
	if overshoot >= 1 {
		overshoot = 0.9999
	}
	if period <= 0 {
		period = 1 - 0.999
	}

	damp := 1.0
	if overshoot > 0 {
		damp = -math.Log(overshoot) / math.Sqrt(math.Pi*math.Pi+math.Log(overshoot))
	}
	// in case of numerical instability in the above, clamp:
	if damp > 1 {
		damp = 1
	}
	if damp < 0 {
		damp = 0
	}
	omega := 1 / period * 2 * math.Pi
	i.damp = damp
	i.omega = omega
}

func (i *interpolator) converged(target, prec float64) bool {
	return math.Abs(i.v) < prec && math.Abs(i.x-target) < prec
}

func (i *interpolator) step(target float64, dt float64) float64 {
	// handle large time gaps (e.g. dropped frames) by taking several steps to not get to
	// much error from the integration. Handle very large time steps (e.g. skipped) by simply
	// setting it to the target, so as to prevent either a very long loop or adding to much
	// energy to the system. It seems a a better fallback to revert to target with 0 velocity
	// than something strange.
	// This should be enough to allow us to use forward euler without stability problems, since
	// we need not be numerically stable but simply look good.

	// TODO: We might actually want a better integrator for faster convergence, such that the
	// entire screen is not updated for longer than is needed. Also some indications that we see
	// artifacts, should should probably use an implicit mehtod?

	maxLength := 2 * i.omega
	maxStep := i.omega / 100

	if dt <= 0 {
		return i.x
	}

	if dt > maxLength {
		i.v = 0
		i.x = target
		return i.x
	}

	tr := dt
	for tr > 0 {
		dt = tr
		if dt > maxStep {
			dt = maxStep
		}
		tr -= dt

		// euler forward, damped harmonic oscillator
		nx := i.x + dt*i.v
		nv := i.v + dt*(-2*i.damp*i.omega*i.v-i.omega*i.omega*(i.x-target))

		i.x = nx
		i.v = nv
	}

	return i.x
}
