package smoothspline

import (
	"fmt"
	"testing"

	"github.com/helloworldpark/gonaturalspline/bspline"
	"github.com/helloworldpark/gonaturalspline/knot"
	"gonum.org/v1/gonum/mat"
)

func TestSmoothSolveRegressionMatrix(t *testing.T) {
	const order = 3
	knots := knot.NewUniformKnot(-10, 0, 10, order)
	fmt.Println(knots, knots.Padding())

	coef := make([]float64, knots.Count()+order)
	simpleSpline := bspline.NewBSplineSimple(order, knots, coef)

	solver := NewSmoothSolver(simpleSpline)
	solver.calcRegressionMatrix()
	fmt.Printf("B: %dx%d \n%0.2v\n", solver.bRegressionMat.RawMatrix().Rows, solver.bRegressionMat.RawMatrix().Cols, mat.Formatted(solver.bRegressionMat))

}
