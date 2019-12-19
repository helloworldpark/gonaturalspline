package bspline

import (
	"fmt"
	"image/color"
	"strconv"
	"testing"

	"github.com/helloworldpark/gonaturalspline/knot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func TestSimpleBSpline(t *testing.T) {
	const order = 3
	knots := knot.NewUniformKnot(0, 1, 11, order)
	fmt.Println(knots)

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = fmt.Sprintf("B-Splines of Order-%d", order)
	p.X.Label.Text = "X"
	ticks := plot.ConstantTicks{}
	for i := 0; i < knots.Count(); i++ {
		v := knots.At(i)
		l := strconv.FormatFloat(v, 'f', 1, 64)
		ticks = append(ticks, plot.Tick{Value: v, Label: l})
	}
	p.X.Tick.Marker = ticks
	p.Y.Label.Text = "Y"

	for m := 0; m <= order; m++ {
		coef := make([]float64, knots.Count()+m)
		simpleSpline := NewBSplineSimple(m, knots, coef)
		f := plotter.NewFunction(simpleSpline.GetBSpline(m).Evaluate)
		f.Samples = 1000
		col := 255.0 * (float64(m) / float64(order))
		f.Color = color.RGBA{R: uint8(col), G: 0, B: 0, A: 255}

		p.Add(f)
	}

	// coef := make([]float64, knots.Count()+order)
	// simpleSpline := NewBSplineSimple(order, knots, coef)
	// simpleSpline.SetCoef(1, 2)
	// // simpleSpline.SetCoef(2, -1.2)
	// simpleSpline.SetCoef(5, -1)
	// // simpleSpline.SetCoef(4, 1)
	// // simpleSpline.SetCoef(5, 0.5)
	// f := plotter.NewFunction(simpleSpline.At)
	// f.Samples = 1000
	// f.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// p.Add(f)

	p.X.Min = knots.At(-1) - 0.5
	p.X.Max = knots.At(knots.Count()+1) + 0.5
	p.Y.Min = -0.5
	p.Y.Max = 3.5

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "/Users/shp/Documents/projects/gonaturalspline/bspline/TestSimpleBSpline.png"); err != nil {
		panic(err)
	}
}
