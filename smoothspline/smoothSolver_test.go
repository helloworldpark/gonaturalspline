package smoothspline

import (
	"fmt"
	"testing"

	"github.com/helloworldpark/gonaturalspline/bspline"
	"github.com/helloworldpark/gonaturalspline/knot"
	"gonum.org/v1/gonum/mat"
)

func TestSmoothSolveRegressionMatrix(t *testing.T) {
	const order = 4
	knots := knot.NewUniformKnot(-10, 0, 11, order)
	fmt.Println(knots, knots.Padding())

	coef := make([]float64, knots.Count()+order)
	simpleSpline := bspline.NewBSplineSimple(order, knots, coef)

	solver := NewSmoothSolver(simpleSpline, 0)
	solver.calcRegressionMatrix()
	solved := solver.RegressionMatrix()
	fmt.Printf("B: %dx%d \n%0.2v\n", solved.RawMatrix().Rows, solved.RawMatrix().Cols, mat.Formatted(solved))
	solver.calcCholesky()
	solved = solver.SolverMatrix()
	fmt.Printf("B: %dx%d \n%0.2v\n", solved.RawMatrix().Rows, solved.RawMatrix().Cols, mat.Formatted(solved))

}
