package knot

import "sort"

import "strings"

import "fmt"

// Knot Definition of Knot
// Knot := {k_0 < k_1 < k_2 < ... < k_count}
// For example, if [0, 1] with interval of 0.1, then the knots are
// Knot = {0, 0.1, 0.2, ... , 0.9, 1.0}
// So, for valid calculation at the end, additional padding should be appended
// i.e. if we need y = spline(1.0), then the knots should be
// Knot = {0, 0.1, 0.2, ... , 0.9, 1.0, 1.1}
// since B-Splines are defined on [k_i, k_(i+1))
type Knot interface {
	// Len
	// Total length of the knots, including paddings
	Len() int

	At(idx int) float64
	Index(x float64) int
	String() string

	IsSorted() bool
	IsUnique() bool
}

type uniformKnot []float64

func NewUniformKnot(start, end float64, count, paddings int) Knot {
	if count <= 0 {
		return nil
	}
	if start >= end {
		return nil
	}
	var knots uniformKnot
	var startIdx = -paddings
	var endIdx = count + paddings
	for i := startIdx; i <= endIdx; i++ {
		knots = append(knots, start+(end-start)*(float64(i)/float64(count)))
	}
	return knots
}

// count + 1
func (k uniformKnot) Len() int {
	return len(k)
}

func (k uniformKnot) IsSorted() bool {
	return sort.Float64sAreSorted(k)
}

func (k uniformKnot) IsUnique() bool {
	if !k.IsSorted() {
		return false
	}
	if k.Len() == 0 {
		return false
	}
	last := k[0]
	for _, f := range k[1:] {
		if last == f {
			return false
		}
		last = f
	}
	return true
}

func (k uniformKnot) At(idx int) float64 {
	if idx < 0 {
		return k[0]
	}
	if idx >= k.Len() {
		return k[k.Len()-1]
	}
	return k[idx]
}

func (k uniformKnot) Index(x float64) int {
	start, end := k[0], k[len(k)-1]
	interval := float64(len(k)-1) / (end - start)
	idx := int(x*interval + start)
	if k.At(idx) <= x && x < k.At(idx+1) {
		return idx
	}
	if idx >= k.Len()-1 {
		return k.Len() - 1
	}
	if idx < 0 {
		return 0
	}
	if x >= k.At(idx+1) {
		for interval < x-k.At(idx) && idx > 0 {
			idx--
		}
		return idx
	}

	for interval < k.At(idx)-x {
		idx++
	}
	return idx
}

func (k uniformKnot) String() string {
	buf := strings.Builder{}
	buf.WriteString("[")
	for i, f := range k {
		buf.WriteString(fmt.Sprintf("%f", f))
		if i < k.Len()-1 {
			buf.WriteString(fmt.Sprintf(", "))
		}
	}
	buf.WriteString("]")
	return buf.String()
}
