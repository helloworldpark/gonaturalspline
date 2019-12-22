package cubicSpline

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

func TestNaturalCubicSpline(t *testing.T) {
	const order = 3
	const lambda = 0.001
	knots := knot.NewUniformKnot(0, 10, 11, order)
	y := []float64{5, 8, 10, 8.5, 4, 0, -3.7, -5, 3.5, -2, 0}

	ncs := NewNaturalCubicSplines(knots, nil)
	ncs.Solve(lambda)
	ncs.Interpolate(y)

	// // Using Solve function
	// var rhs mat.VecDense
	// rhs.MulVec(m.T(), Y)
	// var coef mat.VecDense
	// if err := chol.SolveVecTo(&coef, &rhs); err != nil {
	// 	panic(fmt.Sprintf("Matrix near singular: %+v\n", err))
	// }
	// fmt.Println("Solve L * LT * x = NT * b")
	// var recoveredY mat.VecDense
	// recoveredY.MulVec(&chol, &coef)
	// fmt.Printf("x = %0.4v\nb = %0.4v\n", mat.Formatted(&coef, mat.Prefix("    ")), mat.Formatted(&recoveredY, mat.Prefix("    ")))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = fmt.Sprintf("Natural Cubic Spline")
	p.X.Label.Text = "X"
	ticks := plot.ConstantTicks{}
	for i := 0; i < knots.Count(); i++ {
		v := knots.At(i)
		l := strconv.FormatFloat(v, 'f', 1, 64)
		ticks = append(ticks, plot.Tick{Value: v, Label: l})
	}
	p.X.Tick.Marker = ticks
	p.Y.Label.Text = "Y"

	f := plotter.NewFunction(ncs.At)
	f.Samples = 1000
	f.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}

	p.Add(f)

	data := plotter.XYs{}
	for i := 0; i < len(y); i++ {
		data = append(data, plotter.XY{X: knots.At(i), Y: y[i]})
	}

	scatter, err := plotter.NewScatter(data)
	scatter.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	p.Add(scatter)

	p.X.Min = knots.At(-1) - 0.5
	p.X.Max = knots.At(knots.Count()+1) + 0.5
	p.Y.Min = -12
	p.Y.Max = +12

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "/Users/shp/Documents/projects/gonaturalspline/bspline/TestNaturalCubicSpline.png"); err != nil {
		panic(err)
	}
}
