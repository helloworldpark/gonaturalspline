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
	knots := knot.NewUniformKnot(0, 1, 10, 4)
	const order = 3

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = fmt.Sprintf("B-Splines of Order-%d", order)
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	for i := 0; i < knots.Len(); i++ {
		coef := make([]float64, knots.Len()+order)
		simpleSpline := NewBSplineSimple(order, knots, coef)
		simpleSpline.SetCoef(i+4, 1.0)
		fmt.Printf("B[%d]=%f\n", i, simpleSpline.At(knots.At(i+1)))
		f := plotter.NewFunction(simpleSpline.At)
		f.Samples = 1000

		p.Add(f)
	}

	p.X.Min = -0.5
	p.X.Max = 1.5
	p.Y.Min = 0
	p.Y.Max = 2

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "/Users/shp/Documents/projects/gonaturalspline/bspline/TestSimpleBSpline.png"); err != nil {
		panic(err)
	}
}
