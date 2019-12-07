package knot

import "sort"

import "strings"

import "fmt"

// Knot Definition of Knot
type Knot interface {
	Len() int
	IsSorted() bool
	At(idx int) float64
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
	return knots
}

func (k UniformKnot) Len() int {
	return len(k)
}

func (k UniformKnot) IsSorted() bool {
	return sort.Float64sAreSorted(k)
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
