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
	paddingLeft  map[float64]bool
	paddingRight map[float64]bool
}

// NewArbitraryKnotBuilder Create an arbitrary knot with this
func NewArbitraryKnotBuilder(knots ...float64) *ArbitraryKnotBuilder {
	builder := &ArbitraryKnotBuilder{
		knots:        make(map[float64]bool),
		paddingLeft:  make(map[float64]bool),
		paddingRight: make(map[float64]bool),
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

func (b *ArbitraryKnotBuilder) AppendPaddingLeft(f float64) *ArbitraryKnotBuilder {
	b.paddingLeft[f] = true
	return b
}

func (b *ArbitraryKnotBuilder) AppendPaddingRight(f float64) *ArbitraryKnotBuilder {
	b.paddingRight[f] = true
	return b
}

func (b *ArbitraryKnotBuilder) Build() Knot {
	var knots []float64
	for k := range b.knots {
		knots = append(knots, k)
	}
	sort.Float64s(knots)

	var left []float64
	for k := range b.paddingLeft {
		left = append(left, k)
	}
	sort.Float64s(left)

	var right []float64
	for k := range b.paddingRight {
		right = append(right, k)
	}
	sort.Float64s(right)

	// Knot Check
	if len(knots) <= 1 {
		panic("[Knot] No knots were given")
	}

	// Padding Left check
	if len(left) > 0 && left[len(left)-1] >= knots[0] {
		panic("[Knot] Left padding must be strictly smaller than the main knot")
	}

	// Padding Right Check
	if len(right) > 0 && right[0] <= knots[len(knots)-1] {
		panic("[Knot] Right padding must be strictly bigger than the main knot")
	}

	paddingDiff := len(left) - len(right)
	paddingCount := len(left)
	if paddingDiff < 0 {
		paddingCount = len(right)

		var start = knots[0] - (knots[1] - knots[0])
		if len(left) > 0 {
			start = left[len(left)-1]
		}
		end := knots[0]

		interval := (end - start) / float64(-paddingDiff+2)
		for i := 1; i <= -paddingDiff; i++ {
			left = append(left, start+interval*float64(i))
		}
	} else if paddingDiff > 0 {
		paddingCount = len(left)

		start := knots[len(knots)-1]
		var end = start + (knots[len(knots)-1] - knots[len(knots)-2])
		if len(right) > 0 {
			end = right[0]
		}

		interval := (end - start) / float64(paddingDiff+2)
		for i := 1; i <= paddingDiff; i++ {
			right = append(right, start+interval*float64(i))
		}
	}
	var totalKnots = left
	totalKnots = append(totalKnots, knots...)
	totalKnots = append(totalKnots, right...)
	return &arbitraryKnot{
		knots:   totalKnots,
		padding: paddingCount,
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
