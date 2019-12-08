package bspline

import "github.com/helloworldpark/gonaturalspline/knot"

// BSpline BSpline represents a function that implements B-Spline.
type BSpline interface {
	At(x float64) float64
	AtIdx(idx int) float64
	Knots() knot.Knot
	Order() int
}
