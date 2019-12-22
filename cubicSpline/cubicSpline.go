package cubicSpline

import (
	"github.com/helloworldpark/gonaturalspline/knot"
	"gonum.org/v1/gonum/mat"
)

// CubicSpline Univariate function
type CubicSpline func(float64) float64

// NaturalCubicSplines Reference from:
// p.141-156, T. Hastie et. al., The Elements of Statistical Learning
type NaturalCubicSplines struct {
	splines []CubicSpline
	knots   knot.Knot
	coefs   *mat.VecDense
	lambda  float64

	solverMatrix *mat.Dense
}

// NewNaturalCubicSplines A new pointer of NaturalCubicSpline struct
func NewNaturalCubicSplines(knots knot.Knot, coefs []float64) *NaturalCubicSplines {
	return &NaturalCubicSplines{
		splines: buildNaturalCubicSplines(knots),
		knots:   knots,
		coefs:   mat.NewVecDense(knots.Count(), coefs),
	}
}

// Solve Solve the matrix needed when calculating smoothing spline.
func (ncs *NaturalCubicSplines) Solve(lambda float64) {
	ncs.lambda = lambda
	N := ncs.calcBasisMatrix()
	S := ncs.calcSmoothMatrix()

	var NtN mat.Dense
	NtN.Mul(N.T(), N)

	S.Scale(lambda, S)

	n, _ := NtN.Dims()

	NtNSym := mat.NewSymDense(n, NtN.RawMatrix().Data)
	SSym := mat.NewSymDense(n, S.RawMatrix().Data)
	NtNSym.AddSym(NtNSym, SSym)

	var chol mat.Cholesky
	if ok := chol.Factorize(NtNSym); !ok {
		panic(">>>>>>>>>>")
	}

	var cholInv mat.SymDense
	chol.InverseTo(&cholInv)

	var all mat.Dense
	all.Mul(&cholInv, N.T())

	ncs.solverMatrix = &all
}

// Interpolate Calculate the coefficients interpolating y
func (ncs *NaturalCubicSplines) Interpolate(y []float64) {
	Y := mat.NewVecDense(len(y), y)
	var coefs mat.VecDense
	coefs.MulVec(ncs.solverMatrix, Y)
	ncs.coefs = &coefs
}

// At Calculate the smoothing spline at x
func (ncs *NaturalCubicSplines) At(x float64) float64 {
	var y float64
	for i := 0; i < len(ncs.splines); i++ {
		y += ncs.coefs.AtVec(i) * ncs.splines[i](x)
	}
	return y
}

func (ncs *NaturalCubicSplines) calcBasisMatrix() *mat.Dense {
	n := len(ncs.splines)
	m := mat.NewDense(ncs.knots.Count(), n, nil)
	for i := 0; i < ncs.knots.Count(); i++ {
		x := ncs.knots.At(i)
		for j := 0; j < n; j++ {
			v := ncs.splines[j](x)
			m.Set(i, j, v)
		}
	}
	return m
}

func (ncs *NaturalCubicSplines) calcSmoothMatrix() *mat.Dense {
	n := len(ncs.splines)
	p := mat.NewDense(ncs.knots.Count(), n, nil)

	knotEnd := ncs.knots.At(ncs.knots.Count() - 1)
	knotEndEnd := ncs.knots.At(ncs.knots.Count() - 2)
	for j := 2; j < n; j++ {
		v := 12.0 * (knotEnd - ncs.knots.At(j-2))
		p.Set(j, j, v)
	}
	for j := 2; j < n; j++ {
		knotJ := ncs.knots.At(j - 2)
		diffJ := knotEndEnd - knotJ
		for m := j + 1; m < n; m++ {
			knotM := ncs.knots.At(m - 2)
			diffM := knotEndEnd - knotM
			v := 12.0*(knotEnd-knotEndEnd) + 6.0*(diffM/diffJ)*(2*knotEndEnd-3*knotM+knotJ)
			p.Set(j, m, v)
		}
	}
	for j := 2; j < n; j++ {
		for m := 0; m < j; m++ {
			p.Set(j, m, p.At(m, j))
		}
	}
	return p
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
