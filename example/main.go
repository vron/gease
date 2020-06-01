// Command example illustrates an example easings / animations with gease
package main

import (
	"image/color"
	"log"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/colornames"

	"github.com/vron/gease"
)

const size = 800

var theme *material.Theme

func main() {
	gofont.Register()
	theme = material.NewTheme()
	go func() {
		w := app.NewWindow(app.Size(unit.Px(size), unit.Px(size)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
	ops := &op.Ops{}

	easeHub := buildEasings()

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case pointer.Event:
				updateEasings(e.Position)
				w.Invalidate()
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				ops.Reset()
				gtx := layout.NewContext(ops, e.Queue, e.Config, e.Size)

				t := time.Now()
				if !easeHub.Step(t) {
					op.InvalidateOp{}.Add(ops)
				}

				layoutUI(gtx)
				overlayCircles(ops)

				e.Frame(ops)
			}
		}
	}
}

func layoutUI(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle, Axis: layout.Vertical, Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Flexed(1.0, func(layout.Context) layout.Dimensions {
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			l := material.Label(theme, textSizeEasing.V(), "Gioui easing")
			l.Alignment = text.Middle
			return l.Layout(gtx)
		}),
		layout.Flexed(1.0, func(layout.Context) layout.Dimensions {
			return layout.Dimensions{}
		}),
	)
}

func overlayCircles(ops *op.Ops) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			offx := 50 * (float32(i) - float32(5-1)/2)
			offy := 50 * (float32(j) - float32(5-1)/2)
			p := positionEasings[i*5+j].V()
			p.X += offx
			p.Y += offy
			drawCircle(ops, p, 18, colorEasings[i].V())
		}
	}
}

func drawCircle(ops *op.Ops, c f32.Point, r float32, fill color.RGBA) {
	stack := op.StackOp{}
	stack.Push(ops)

	re := f32.Rect(c.X-r, c.Y-r, c.X+r, c.Y+r)
	clip.Rect{
		Rect: re,
		SE:   r,
		SW:   r,
		NW:   r,
		NE:   r,
	}.Add(ops)
	paint.ColorOp{Color: fill}.Add(ops)
	paint.PaintOp{Rect: re}.Add(ops)

	stack.Pop()
}

var positionEasings []*gease.PointEasing
var colorEasings []*gease.ColorEasing
var textSizeEasing *gease.UnitEasing

func buildEasings() *gease.Hub {
	periods := []time.Duration{
		500 * time.Millisecond,
		750 * time.Millisecond,
		1000 * time.Millisecond,
		1500 * time.Millisecond,
		2000 * time.Millisecond,
	}
	overshoots := []float64{0, 0.1, 0.15, 0.2, 0.25}

	positionEasings = []*gease.PointEasing{}
	colorEasings = []*gease.ColorEasing{}
	for _, p := range periods {
		for _, o := range overshoots {
			positionEasings = append(positionEasings, gease.Point(f32.Point{X: 0, Y: 0}).Configure(o, p))
		}
		colorEasings = append(colorEasings, gease.Color(colornames.Red100).Configure(p*2))
	}

	textSizeEasing = gease.Unit(unit.Sp(20))

	hub := gease.NewHub()
	hub.Add(textSizeEasing)
	for _, e := range positionEasings {
		hub.Add(e)
	}
	for _, e := range colorEasings {
		hub.Add(e)
	}
	return hub
}

func updateEasings(pos f32.Point) {
	for _, es := range positionEasings {
		es.Target(pos)
	}
	for _, es := range colorEasings {
		alpha := float64(pos.X*pos.Y) / float64(size*size)
		alpha = 0.0
		if pos.X > pos.Y {
			col := colornames.Red400
			col.A = 255 - uint8(alpha*255+0.5)
			col.R = uint8(uint16(col.R) * uint16(col.A) / 255)
			col.G = uint8(uint16(col.G) * uint16(col.A) / 255)
			col.B = uint8(uint16(col.B) * uint16(col.A) / 255)
			es.Target(col)
		} else {

			col := colornames.Green400
			col.A = 255 - uint8(alpha*255+0.5)
			col.R = uint8(uint16(col.R) * uint16(col.A) / 255)
			col.G = uint8(uint16(col.G) * uint16(col.A) / 255)
			col.B = uint8(uint16(col.B) * uint16(col.A) / 255)
			es.Target(col)
		}
	}
	textSizeEasing.Target(unit.Sp(pos.Y/size*40 + 15))
}
