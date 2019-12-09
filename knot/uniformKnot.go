package knot

import "sort"

import "strings"

import "fmt"

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
	for i := 0; i < paddings; i++ {
		knots.knots = append(knots.knots, start)
	}
	for i := 0; i < count; i++ {
		knots.knots = append(knots.knots, start+(end-start)*(float64(i)/float64(count-1)))
	}
	for i := 0; i < paddings; i++ {
		knots.knots = append(knots.knots, end)
	}
	knots.padding = paddings
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
	idx += k.Padding()
	if idx < 0 {
		return k.knots[0]
	}
	if idx >= k.Len() {
		return k.knots[k.Len()-1]
	}
	return k.knots[idx]
}

func (k *uniformKnot) Index(x float64) int {
	idx := sort.Search(len(k.knots), func(i int) bool {
		return k.knots[i] >= x
	})
	if idx == 0 {
		return -k.Padding()
	}
	// -1: since sort.Search returns smallest idx s.t. k.knots[idx] >= x,
	//     the knot should be knot_(idx-1) <= x < knot_idx
	return idx - 1 - k.Padding()
}

func (k *uniformKnot) String() string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("UniformKnot(Count: %d, Padding: %d)[", k.Count(), k.Padding()))
	for i, f := range k.knots {
		buf.WriteString(fmt.Sprintf("%f", f))
		if i < k.Len()-1 {
			buf.WriteString(fmt.Sprintf(", "))
		}
	}
	buf.WriteString("]")
	return buf.String()
}
