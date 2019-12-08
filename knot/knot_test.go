package knot

import (
	"fmt"
	"testing"
)

func TestUniformKnot(t *testing.T) {
	knot := NewUniformKnot(0, 10, 100, 4)
	if !knot.IsSorted() {
		t.Fatal("[Knot] Knot is not sorted!")
	}
	if !knot.IsUnique() {
		t.Fatal("[Knot] Knot is not unique!")
	}
	t.Logf("[Knot] Length = %d\n", knot.Len())
	for i := -10; i <= 110; i++ {
		t.Logf("[Knot] Knot[%d] = %f\n", i, knot.At(i))
	}
	t.Logf("[Knot] Knot = %v\n", knot)
}

func TestKnotIndex(t *testing.T) {
	knots := NewUniformKnot(2, 3, 10, 4)
	for x := 0.0; x < 2.0; x += 0.00001 {
		fmt.Printf("x[%d]=%f\n", knots.Index(x), x)
	}
}
