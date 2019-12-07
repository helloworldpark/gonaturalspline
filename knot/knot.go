package knot

import "sort"

import "strings"

import "fmt"

// Knot Definition of Knot
type Knot interface {
	Len() int
	IsSorted() bool
	IsUnique() bool
	At(idx int) float64
	Index(x float64) int
	String() string
}

type UniformKnot []float64

func NewUniformKnot(start, end float64, count int) Knot {
	if count <= 0 {
		return nil
	}
	if start >= end {
		return nil
	}
	var knots UniformKnot
	for i := 0; i < count; i++ {
		knots = append(knots, (end-start)*(float64(i)/float64(count)))
	}
	knots = append(knots, end)
	return knots
}

// count + 1
func (k UniformKnot) Len() int {
	return len(k)
}

func (k UniformKnot) IsSorted() bool {
	return sort.Float64sAreSorted(k)
}

func (k UniformKnot) IsUnique() bool {
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

func (k UniformKnot) At(idx int) float64 {
	if idx < 0 {
		return k[0]
	}
	if idx >= k.Len() {
		return k[k.Len()-1]
	}
	return k[idx]
}

func (k UniformKnot) Index(x float64) int {
	start, end := k[0], k[len(k)-1]
	interval := (end - start) / float64(len(k)-1)
	idx := int(x*interval + start)
	if k.At(idx) <= x && x < k.At(idx+1) {
		return idx
	}
	if x >= k.At(idx+1) {
		for interval < x-k.At(idx) && idx < len(k) {
			idx++
		}
		return idx
	}

	for interval > k.At(idx)-x && x > 0 {
		idx--
	}
	return idx
}

func (k UniformKnot) String() string {
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
