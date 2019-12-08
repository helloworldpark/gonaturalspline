package smoothspline

import (
	"fmt"

	"github.com/helloworldpark/gonaturalspline/bspline"
	"gonum.org/v1/gonum/mat"
)

type SmoothSolver struct {
	bSpline          bspline.BSpline
	bRegressionMat   *mat.Dense
	bCholeskyedLower *mat.Dense
	bPenaltyMat      *mat.Dense
}

func NewSmoothSolver(spline bspline.BSpline) *SmoothSolver {
	return &SmoothSolver{
		bSpline: spline,
	}
}

func (solver *SmoothSolver) calcRegressionMatrix() {
	order := solver.bSpline.Order()
	N := solver.bSpline.Knots().Count()

	B := mat.NewDense(N, N+order, nil)
	for i := 0; i < N; i++ {
		for j := 0; j < N+order; j++ {
			x := solver.bSpline.Knots().At(i)
			fmt.Println(i, x)
			v := solver.bSpline.GetBSpline(j)(x)
			B.Set(i, j, v)
		}
	}
	solver.bRegressionMat = B
}
