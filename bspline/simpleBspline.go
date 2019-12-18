package bspline

import (
	"github.com/helloworldpark/gonaturalspline/knot"
)

type bFunc func(float64) float64
type bSplineSimple struct {
	knots    knot.Knot
	order    int
	bsplines []bFunc
	coefs    []float64
}

func NewBSplineSimple(order int, knot knot.Knot, coef []float64) BSpline {
	// TODO:
	// assert: order >= 1, knot.Count() > 0, knot.Count() + order == len(coef)
	bsplines := constructBSplines(order, knot)
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
		v += b.GetCoef(idx+m) * b.GetBSpline(idx+m)(x)
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
func (b *bSplineSimple) GetBSpline(idx int) bFunc {
	idx += b.order
	if idx < 0 {
		return b.bsplines[0]
	}
	if idx >= len(b.bsplines) {
		return b.bsplines[len(b.bsplines)-1]
	}
	return b.bsplines[idx]
}

func constructZeroFunction() bFunc {
	return func(float64) float64 {
		return 0
	}
}

func constructHaarFunction(k1, k2 float64) bFunc {
	return func(x float64) float64 {
		if k1 <= x && x < k2 {
			return 1
		}
		return 0
	}
}

func constructBsplinesRecursively(a1, a2, t1, t2 float64, fa, ft bFunc) bFunc {
	if a1 == a2 {
		if t1 == t2 {
			return constructZeroFunction()
		}
		// No a1a2 coefficient
		return func(x float64) float64 {
			return (t2 - x) / (t2 - t1) * ft(x)
		}
	} else if t1 == t2 {
		// No t1t2 coefficient
		return func(x float64) float64 {
			return (x - a1) / (a2 - a1) * fa(x)
		}
	}
	// No zero coefficient
	return func(x float64) float64 {
		return (x-a1)/(a2-a1)*fa(x) + (t2-x)/(t2-t1)*ft(x)
	}
}

// [i]bFunc: {order}-order B-Spline associated with j-th knot B_j^i(x)
func constructBSplines(order int, knots knot.Knot) []bFunc {
	if order < 0 {
		return nil
	}
	var splines []bFunc
	// Order 0
	for idx := -order; idx < knots.Count()+order; idx++ {
		k1, k2 := knots.At(idx), knots.At(idx+1)
		var fcn bFunc
		if k1 == k2 {
			fcn = constructZeroFunction()
		} else {
			fcn = constructHaarFunction(k1, k2)
		}
		splines = append(splines, fcn)
	}

	// Order 1~ : Recursive
	for m := 1; m <= order; m++ {
		for idx := -order; idx < knots.Count()+order-m; idx++ {
			a1, a2 := knots.At(idx), knots.At(idx+m)
			t1, t2 := knots.At(idx+1), knots.At(idx+m+1)
			fa := splines[idx+order]
			ft := splines[idx+1+order]
			fcn := constructBsplinesRecursively(a1, a2, t1, t2, fa, ft)
			splines[idx+order] = fcn
		}
	}
	splines = splines[:knots.Count()+order]
	return splines
}

type bSplineFunc interface {
}

type bSplineOrder struct {
}
