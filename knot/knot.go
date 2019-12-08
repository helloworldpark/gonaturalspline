package knot

import "sort"

import "strings"

import "fmt"

// Knot Definition of Knot
//     Knot := {k_0 < k_1 < k_2 < ... < k_count}
// For example, if [0, 1] with interval of 0.1, then the knots are
//     Knot = {0, 0.1, 0.2, ... , 0.9, 1.0}
// So, for valid calculation at the both ends, additional padding should be appended
// i.e. if we need y = spline(1.0), then the knots should be
//     Knot = {0, 0.1, 0.2, ... , 0.9, 1.0, 1.1}
//     since B-Splines are defined on [k_i, k_(i+1)).
// If padding is included, then we interpret knots as:
//     k_-p, k_(-p+1), ... , k_-1, k_0, k_1, ... , k_count, k_(count+1), ... , k_(count+p)
//     --------------------------  ^^^^^^^^^^^^^^^^^^^^^^^  ------------------------------
//             PADDINGS                    KNOTS                       PADDINGS
type Knot interface {
	// Len
	// Total length of the knots, including paddings
	Len() int
	// Padding
	// How many paddings are on each end?
	Padding() int
	// Count
	// Length of the knots without paddings
	Count() int

	// Will return value considering padding, i.e. if padding = 4, then At(0) = Knot[4]
	At(idx int) float64
	// Will return value considering padding, i.e. if padding = 4, then Index(0.0) = 0
	Index(x float64) int
	String() string

	IsSorted() bool
	IsUnique() bool
}

type uniformKnot struct {
	knots   []float64
	padding int
}

// NewUniformKnot Creates a new Knot with uniform intervals
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
		knots.knots = append(knots.knots, start+(end-start)*(float64(i)/float64(count)))
	}
	return &knots
}

func (k *uniformKnot) Len() int {
	return len(k.knots)
}

func (k *uniformKnot) Padding() int {
	return k.padding
}

func (k *uniformKnot) Count() int {
	return len(k.knots) - 2*k.padding
}

func (k *uniformKnot) IsSorted() bool {
	return sort.Float64sAreSorted(k.knots)
}

func (k *uniformKnot) IsUnique() bool {
	if !k.IsSorted() {
		return false
	}
	if k.Len() == 0 {
		return false
	}
	last := k.knots[0]
	for _, f := range k.knots[1:] {
		if last == f {
			return false
		}
		last = f
	}
	return true
}

func (k *uniformKnot) At(idx int) float64 {
	if idx+k.Padding() < 0 {
		return k.knots[0]
	}
	if idx >= k.Count()+k.Padding() {
		return k.knots[k.Len()-1]
	}
	return k.knots[idx+k.Padding()]
}

func (k *uniformKnot) Index(x float64) int {
	start, end := k.knots[0], k.knots[len(k.knots)-1]
	interval := float64(len(k.knots)-1) / (end - start)
	idx := int(x*interval + start)
	if k.At(idx) <= x && x < k.At(idx+1) {
		return idx - k.Padding()
	}
	if idx >= k.Len()-1-k.Padding() {
		return k.Len() - 1 - k.Padding()
	}
	if idx < 0 {
		return k.Padding()
	}
	if x >= k.At(idx+1) {
		for interval < x-k.At(idx) && idx > 0 {
			idx--
		}
		return idx - k.Padding()
	}

	for interval < k.At(idx)-x {
		idx++
	}
	return idx - k.Padding()
}

func (k *uniformKnot) String() string {
	buf := strings.Builder{}
	buf.WriteString("[")
	for i, f := range k.knots {
		buf.WriteString(fmt.Sprintf("%f", f))
		if i < k.Len()-1 {
			buf.WriteString(fmt.Sprintf(", "))
		}
	}
	buf.WriteString("]")
	return buf.String()
}
