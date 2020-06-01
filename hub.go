package gease

import (
	"reflect"
	"time"
)

// Easing is the interface an easing must satisfy to be used in a Hub.
type Easing interface {
	// SetTime is used to advance the easing in time
	Step(t time.Time) (converged bool)
}

// A Hub is a convenience wrapper for cases with constantly live easings.
type Hub struct {
	easings map[uintptr]Easing
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		easings: make(map[uintptr]Easing, 100),
	}
}

// Add adds one or more easing to the hub.
func (h *Hub) Add(ess ...Easing) {
	if ess == nil {
		return
	}
	for _, es := range ess {
		// TODO: Ensure the below is actually safe, such that the same object will always have the same
		// pointer and not be moved? If it was e.g. moved by the GC we would risk caliing .Step() at the
		// same easing twice - which would slow down performance a bit but have no other bad consequences.
		adr := reflect.ValueOf(es).Pointer()
		h.easings[adr] = es
	}
}

// Remove the easing from the hub, typically used such that it can be garbage
// collected if it will not be used further.
func (h *Hub) Remove(es Easing) {
	adr := reflect.ValueOf(es).Pointer()
	delete(h.easings, adr)
}

// Step all easings in the Hub, if no time is provided time.Now() will be used and
// if more than one time is provided only the first one will be used. If all underlying
// easings have converged true will be returned, else false is returned.
func (h *Hub) Step(t ...time.Time) (converged bool) {
	converged = true
	var tt time.Time
	if t != nil && len(t) > 0 {
		tt = t[0]
	}
	for _, es := range h.easings {
		if !es.Step(tt) {
			converged = false
		}
	}
	return converged
}
