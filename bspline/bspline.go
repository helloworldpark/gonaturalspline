package bspline

import "github.com/helloworldpark/gonaturalspline/knot"

type BSplineGroup interface {
	At(x float64) float64
	AtIdx(idx int) float64
	Knots() knot.Knot
	Order() int
}
