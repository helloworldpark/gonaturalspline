package cubicSpline

import (
	"github.com/helloworldpark/gonaturalspline/knot"
	"gonum.org/v1/gonum/mat"
)

type CubicSpline func(float64) float64

type NaturalCubicSplines struct {
	Splines []CubicSpline
	Knots   knot.Knot
}

func NewNaturalCubicSplines(knots knot.Knot) *NaturalCubicSplines {
	return &NaturalCubicSplines{
		Splines: buildNaturalCubicSplines(knots),
		Knots:   knots,
	}
}

func (ncs *NaturalCubicSplines) At(x float64, coefs *mat.VecDense) float64 {
	var y float64
	for i := 0; i < len(ncs.Splines); i++ {
		y += coefs.AtVec(i) * ncs.Splines[i](x)
	}
	return y
}

func (ncs *NaturalCubicSplines) Matrix() *mat.Dense {
	n := len(ncs.Splines)
	m := mat.NewDense(ncs.Knots.Count(), n, nil)
	for i := 0; i < ncs.Knots.Count(); i++ {
		x := ncs.Knots.At(i)
		for j := 0; j < n; j++ {
			v := ncs.Splines[j](x)
			m.Set(i, j, v)
		}
	}
	return m
}

func piecewiseCubic(k float64) CubicSpline {
	return func(x float64) float64 {
		if x < k {
			return 0.0
		}
		t := x - k
		return t * t * t
	}
}

func buildNaturalCubicSplines(knots knot.Knot) []CubicSpline {
	splines := make([]CubicSpline, knots.Count())
	splines[0] = func(float64) float64 { return 1 }
	splines[1] = func(x float64) float64 { return x }

	knotEnd := knots.At(knots.Count() - 1)
	pEnd := piecewiseCubic(knotEnd)
	dEnd := func(x float64) float64 {
		knotLastToSecond := knots.At(knots.Count() - 2)
		p := piecewiseCubic(knotLastToSecond)
		return (p(x) - pEnd(x)) / (knotEnd - knotLastToSecond)
	}

	for k := 0; k < knots.Count()-2; k++ {
		l := knots.At(k)
		splines[k+2] = func(x float64) float64 {
			p := piecewiseCubic(l)
			return (p(x)-pEnd(x))/(knotEnd-l) - dEnd(x)
		}
	}
	return splines
}
