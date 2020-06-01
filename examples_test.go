package gease

import (
	"fmt"
	"time"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

func ExampleColor() {
	red, green := colornames.Red400, colornames.Green400

	es := Color(red)
	es.Target(green)

	tt := time.Now()
	for !es.Step(tt) {
		// one would usually do: tt = time.Now()
		tt = tt.Add(time.Millisecond * 20)

		// use the color for drawing
		col := es.V()
		_ = col
	}

	// once the animation is converged we have the expected value
	fmt.Println(es.V() == green)
	// Output: true
}

func ExampleUnit() {

}
