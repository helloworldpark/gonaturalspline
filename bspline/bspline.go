package bspline

import "github.com/helloworldpark/gonaturalspline/knot"

// BSpline BSpline represents a function that implements B-Spline.
type BSpline interface {
	At(x float64) float64
	Knots() knot.Knot
	Order() int
	SetCoef(idx int, v float64)
	GetCoef(idx int) float64
	GetBSpline(idx int) bFunc
}
