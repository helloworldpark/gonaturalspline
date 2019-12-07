package bspline

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func Hello() {
	// Construct a symmetric positive definite matrix.
	zero := mat.NewDense(100, 100, nil)
	fmt.Printf("Zero = %v\n", zero.Trace())
}
