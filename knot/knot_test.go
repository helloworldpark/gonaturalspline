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

func TestArbitraryKnot(t *testing.T) {
	builder := NewArbitraryKnotBuilder()
	builder = builder.Append(0.0).Append(1.0).Append(1.5).Append(0.1).Append(3.1415)
	builder = builder.AppendPaddingLeft(-1.0).AppendPaddingLeft(-2.0)
	builder = builder.AppendPaddingRight(3.15).AppendPaddingRight(3.155)
	knot := builder.Build()

	fmt.Println(knot)
}
