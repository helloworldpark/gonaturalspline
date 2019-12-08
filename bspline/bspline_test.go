package bspline

import (
	"fmt"
	"testing"

	"github.com/helloworldpark/gonaturalspline/knot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func TestSimpleBSpline(t *testing.T) {
	const order = 3
	knots := knot.NewUniformKnot(0, 1, 10, order)

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = fmt.Sprintf("B-Splines of Order-%d", order)
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	for i := 0; i < knots.Count(); i++ {
		coef := make([]float64, knots.Count()+order)
		simpleSpline := NewBSplineSimple(order, knots, coef)
		simpleSpline.SetCoef(i+order, 1.0)
		x := knots.At(i + order)
		fmt.Printf("B[%d](%f)=%f\n", i, x, simpleSpline.At(x))
		f := plotter.NewFunction(simpleSpline.At)
		f.Samples = 1000

		p.Add(f)
	}

	p.X.Min = -0.2
	p.X.Max = 1.2
	p.Y.Min = 0
	p.Y.Max = 2

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "/Users/shp/Documents/projects/gonaturalspline/bspline/TestSimpleBSpline.png"); err != nil {
		panic(err)
	}
}
