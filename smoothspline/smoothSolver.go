package smoothspline

import (
	"fmt"

	"github.com/helloworldpark/gonaturalspline/bspline"
	"gonum.org/v1/gonum/mat"
)

type SmoothSolver struct {
	bSpline        bspline.BSpline
	bRegressionMat *mat.Dense
	bSolvedMat     *mat.Dense
	bPenaltyMat    *mat.Dense // scale by lambda at calculation
	lambda         float64
}

func NewSmoothSolver(spline bspline.BSpline, lambda float64) *SmoothSolver {
	return &SmoothSolver{
		bSpline: spline,
		lambda:  lambda,
	}
}

func (solver *SmoothSolver) calcRegressionMatrix() {
	order := solver.bSpline.Order()
	N := solver.bSpline.Knots().Count()

	B := mat.NewDense(N, N+order+1, nil)
	for i := 0; i < N; i++ {
		x := solver.bSpline.Knots().At(i)
		for j := 0; j < N+order; j++ {
			v := solver.bSpline.GetBSpline(j)(x)
			B.Set(i, j, v)
		}
	}
	solver.bRegressionMat = B
}

func (solver *SmoothSolver) RegressionMatrix() *mat.Dense {
	if solver.bRegressionMat == nil {
		return nil
	}
	regMat := mat.NewDense(solver.bRegressionMat.RawMatrix().Rows, solver.bRegressionMat.RawMatrix().Cols, nil)
	_, _ = regMat.Copy(solver.bRegressionMat)
	return regMat
}

func (solver *SmoothSolver) calcCholesky() {
	regMat := solver.bRegressionMat
	cols := regMat.RawMatrix().Cols
	btb := mat.NewDense(cols, cols, nil)
	btb.Mul(regMat.T(), regMat)
	if solver.bPenaltyMat != nil {
		btb.Add(btb, solver.bPenaltyMat)
	}
	btbSym := mat.NewSymDense(cols, btb.RawMatrix().Data)
	r, c := btbSym.Dims()
	fmt.Printf("BTB: %dx%d \n%0.2v\n", r, c, mat.Formatted(btbSym))

	var chol mat.Cholesky
	if ok := chol.Factorize(btbSym); !ok {
		panic(">>>>>>>>>>")
	}
	L := mat.NewTriDense(cols, mat.Lower, nil)
	chol.LTo(L)

	L.InverseTri(L)
	btb.Mul(L, regMat.T())
	L.InverseTri(L)
	btb.Mul(L, btb)

	solver.bSolvedMat = btb
}

func (solver *SmoothSolver) SolverMatrix() *mat.Dense {
	if solver.bSolvedMat == nil {
		return nil
	}
	solMat := mat.NewDense(solver.bSolvedMat.RawMatrix().Rows, solver.bSolvedMat.RawMatrix().Cols, nil)
	_, _ = solMat.Copy(solver.bSolvedMat)
	return solMat
}
