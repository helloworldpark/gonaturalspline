package bspline

import (
	"github.com/helloworldpark/gonaturalspline/knot"
)

type bSplineSimple struct {
	knots    knot.Knot
	order    int
	bsplines []BSplineFunc
	coefs    []float64
}

func NewBSplineSimple(order int, knot knot.Knot, coef []float64) BSpline {
	// TODO:
	// assert: order >= 1, knot.Count() > 0, knot.Count() + order == len(coef)
	bsplines := buildBSplines(order, knot)
	return &bSplineSimple{
		knots:    knot,
		order:    order,
		bsplines: bsplines,
		coefs:    coef,
	}
}

func (b *bSplineSimple) At(x float64) float64 {
	idx := b.knots.Index(x)
	var v float64
	for m := -b.order; m <= 0; m++ {
		v += b.GetCoef(idx+m) * b.GetBSpline(idx+m).Evaluate(x)
	}
	return v
}

func (b *bSplineSimple) Knots() knot.Knot {
	return b.knots
}

func (b *bSplineSimple) Order() int {
	return b.order
}

func (b *bSplineSimple) SetCoef(idx int, v float64) {
	if 0 <= idx && idx < len(b.coefs) {
		b.coefs[idx] = v
	} else {
		// TODO: Warn
	}
}

func (b *bSplineSimple) GetCoef(idx int) float64 {
	idx += b.order
	if idx < 0 {
		return 0.0
	}
	if idx >= len(b.coefs) {
		return 0.0
	}
	return b.coefs[idx]
}

// GetBSpline find cached B-Spline function from index of the knot
func (b *bSplineSimple) GetBSpline(idx int) BSplineFunc {
	// idx += b.order
	if idx < 0 {
		return b.bsplines[0]
	}
	if idx >= len(b.bsplines) {
		return b.bsplines[len(b.bsplines)-1]
	}
	return b.bsplines[idx]
}

///////////////////////////////////////

type BSplineFunc interface {
	Evaluate(float64) float64
}

type bSplineHaar struct {
	knots [2]float64
}

func (b *bSplineHaar) Evaluate(x float64) float64 {
	if b.knots[0] <= x && x < b.knots[1] {
		return 1
	}
	return 0
}

type bSplineOrder struct {
	leftRamp   BSplineFunc
	rightRamp  BSplineFunc
	leftKnots  [2]float64 // ti, ti+o
	rightKnots [2]float64 // ti+1, ti+o+1
	order      int
}

func (b *bSplineOrder) Evaluate(x float64) float64 {
	var y float64
	if b.leftKnots[0] != b.leftKnots[1] {
		y += ((x - b.leftKnots[0]) / (b.leftKnots[1] - b.leftKnots[0])) * b.leftRamp.Evaluate(x)
	}
	if b.rightKnots[0] != b.rightKnots[1] {
		y += ((b.rightKnots[1] - x) / (b.rightKnots[1] - b.rightKnots[0])) * b.rightRamp.Evaluate(x)
	}
	return y
}

func newBSplineOrder(t1, t2 float64, t3, t4 float64, order int, lspline, rspline BSplineFunc) *bSplineOrder {
	return &bSplineOrder{
		leftRamp:   lspline,
		rightRamp:  rspline,
		order:      order,
		leftKnots:  [2]float64{t1, t2},
		rightKnots: [2]float64{t3, t4},
	}
}

func buildBSplines(order int, knots knot.Knot) []BSplineFunc {
	if order < 0 {
		return nil
	}
	var splines []BSplineFunc
	// Order 0
	for idx := -order; idx < knots.Count()+order; idx++ {
		k1, k2 := knots.At(idx), knots.At(idx+1)
		splines = append(splines, &bSplineHaar{knots: [2]float64{k1, k2}})
	}

	// Order 1~ : Recursive
	for m := 1; m <= order; m++ {
		for idx := -order; idx < knots.Count()+order-m; idx++ {
			a1, a2 := knots.At(idx), knots.At(idx+m)
			t1, t2 := knots.At(idx+1), knots.At(idx+m+1)
			fa := splines[idx+order]
			ft := splines[idx+1+order]
			fcn := newBSplineOrder(a1, a2, t1, t2, m, fa, ft)
			splines[idx+order] = fcn
		}
	}
	splines = splines[:knots.Count()+order]
	return splines
}
