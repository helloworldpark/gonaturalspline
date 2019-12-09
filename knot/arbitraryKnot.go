package knot

import "sort"

import "strings"

import "fmt"

type arbitraryKnot struct {
	knots   []float64
	padding int
}

type ArbitraryKnotBuilder struct {
	knots        map[float64]bool
	paddingCount int
}

// NewArbitraryKnotBuilder Create an arbitrary knot with this
func NewArbitraryKnotBuilder(paddingCount int, knots ...float64) *ArbitraryKnotBuilder {
	builder := &ArbitraryKnotBuilder{
		knots:        make(map[float64]bool),
		paddingCount: paddingCount,
	}
	for _, f := range knots {
		builder.knots[f] = true
	}
	return builder
}

func (b *ArbitraryKnotBuilder) Append(f float64) *ArbitraryKnotBuilder {
	b.knots[f] = true
	return b
}

func (b *ArbitraryKnotBuilder) Build() Knot {
	var knots []float64
	for k := range b.knots {
		knots = append(knots, k)
	}
	sort.Float64s(knots)

	// Knot Check
	if len(knots) <= 1 {
		panic("[Knot] No knots were given")
	}

	var totalKnots []float64
	for i := 0; i < b.paddingCount; i++ {
		totalKnots = append(totalKnots, knots[0])
	}
	totalKnots = append(totalKnots, knots...)
	for i := 0; i < b.paddingCount; i++ {
		totalKnots = append(totalKnots, knots[len(knots)-1])
	}
	return &arbitraryKnot{
		knots:   totalKnots,
		padding: b.paddingCount,
	}
}

func (k *arbitraryKnot) Len() int {
	return len(k.knots)
}

func (k *arbitraryKnot) Padding() int {
	return k.padding
}

func (k *arbitraryKnot) Count() int {
	return len(k.knots) - 2*k.padding
}

func (k *arbitraryKnot) IsSorted() bool {
	return sort.Float64sAreSorted(k.knots)
}

func (k *arbitraryKnot) IsUnique() bool {
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

func (k *arbitraryKnot) At(idx int) float64 {
	if idx+k.Padding() < 0 {
		return k.knots[0]
	}
	if idx >= k.Count()+k.Padding() {
		return k.knots[k.Len()-1]
	}
	return k.knots[idx+k.Padding()]
}

func (k *arbitraryKnot) Index(x float64) int {
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

func (k *arbitraryKnot) String() string {
	buf := strings.Builder{}
	buf.WriteString("ArbitraryKnot[")
	for i, f := range k.knots {
		buf.WriteString(fmt.Sprintf("%f", f))
		if i < k.Len()-1 {
			buf.WriteString(fmt.Sprintf(", "))
		}
	}
	buf.WriteString("]")
	return buf.String()
}
