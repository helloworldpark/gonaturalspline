package knot

import "testing"

func TestUniformKnot(t *testing.T) {
	knot := NewUniformKnot(0, 10, 100)
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
