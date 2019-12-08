package knot

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
