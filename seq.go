package seq

import (
	"cmp"
	"errors"
	"iter"

	"golang.org/x/exp/constraints" //nolint:exptostd // Integer and Float constraints are not available in std
)

// Aggregate applies an accumulator function over a sequence.
func Aggregate[V, A any](seq iter.Seq[V], init A, f func(A, V) A) A {
	acc := init
	for v := range seq {
		acc = f(acc, v)
	}

	return acc
}

// All determines if all values of a sequence satisfy a condition.
//
// This returns true if the sequence was empty.
func All[V any](seq iter.Seq[V], f func(V) bool) bool {
	for v := range seq {
		if !f(v) {
			return false
		}
	}

	return true
}

// Any determines if a sequence has any values.
func Any[V any](seq iter.Seq[V]) bool {
	for range seq {
		return true
	}

	return false
}

// AnyFunc determines if a sequence has any value that satisfies a predicate.
func AnyFunc[V any](seq iter.Seq[V], f func(V) bool) bool {
	for v := range seq {
		if f(v) {
			return true
		}
	}

	return false
}

// Append adds values to the end of a sequence.
func Append[V any](seq iter.Seq[V], vals ...V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !yield(v) {
				return
			}
		}

		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
	}
}

// Average computes the average of the values in a sequence.
func Average[V constraints.Integer | constraints.Float](seq iter.Seq[V]) float64 {
	var sum V
	var count int

	for v := range seq {
		sum += v
		count++
	}

	if count == 0 {
		return 0
	}

	return float64(sum) / float64(count)
}

// Chunk splits the values of a sequence into slices of a given size at most.
func Chunk[V any](seq iter.Seq[V], size int) iter.Seq[[]V] {
	return func(yield func([]V) bool) {
		var chunk []V

		for v := range seq {
			if chunk == nil {
				// Lazily allocate an array for the chunk
				chunk = make([]V, 0, size)
			}

			chunk = append(chunk, v)
			if len(chunk) == size {
				if !yield(chunk) {
					return
				}

				// Reset the chunk; a new one will be allocated if there are more values
				chunk = nil
			}
		}

		// Make sure to return a partial chunk
		if len(chunk) > 0 {
			yield(chunk)
		}
	}
}

// Concat concatenates multiples sequences into a single sequence.
func Concat[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Contains determines if a sequence contains a value.
func Contains[V comparable](seq iter.Seq[V], val V) bool {
	for v := range seq {
		if v == val {
			return true
		}
	}

	return false
}

// Count returns the number of values in a sequence.
//
// This iterates over the entire sequence to count the values.
// Use [Any] instead if you only need to check whether the sequence has any values.
// Use [Single] instead if you need to check whether the sequence has exactly one value.
func Count[V any](seq iter.Seq[V]) int {
	count := 0
	for range seq {
		count++
	}

	return count
}

// CountFunc returns the number of values in a sequence that satisfy a predicate.
//
// This iterates over the entire sequence to count the values.
// Use [AnyFunc] instead if you only need to check whether the sequence has any values that satisfy
// a predicate.
// Use [SingleFunc] instead if you need to check whether the sequence has exactly one value that
// satisfies a predicate.
func CountFunc[V any](seq iter.Seq[V], f func(V) bool) int {
	count := 0

	for v := range seq {
		if f(v) {
			count++
		}
	}

	return count
}

// Empty returns an empty sequence.
func Empty[V any]() iter.Seq[V] {
	return func(func(V) bool) {}
}

// Equal determines if two sequences are equal.
func Equal[V comparable](seq1, seq2 iter.Seq[V]) bool {
	nextV2, done := iter.Pull(seq2)
	defer done()

	for v1 := range seq1 {
		v2, ok := nextV2()
		if !ok {
			// seq2 has fewer values
			return false
		}

		if v1 != v2 {
			// each element must be equal
			return false
		}
	}

	// not equal if seq2 has more values
	_, more := nextV2()

	return !more
}

// EqualFunc determines if two sequences are equal using a function to compare values.
func EqualFunc[V any](seq1, seq2 iter.Seq[V], f func(V, V) bool) bool {
	nextV2, done := iter.Pull(seq2)
	defer done()

	for v1 := range seq1 {
		v2, ok := nextV2()
		if !ok {
			// seq2 has fewer values
			return false
		}

		if !f(v1, v2) {
			// check if each value is equal
			return false
		}
	}

	// not equal if seq2 has more values
	_, more := nextV2()

	return !more
}

// First returns the first value of a sequence.
//
// A second return value indicates whether the sequence contained any values.
func First[V any](seq iter.Seq[V]) (V, bool) {
	var v V
	for v = range seq {
		return v, true
	}

	return v, false
}

// FirstFunc returns the first value of a sequence that satisfies a predicate.
//
// A second return value indicates whether the sequence contained any value that satisfied the predicate.
func FirstFunc[V any](seq iter.Seq[V], f func(V) bool) (V, bool) {
	var zero V

	for v := range seq {
		if f(v) {
			return v, true
		}
	}

	return zero, false
}

// Last returns the last value of a sequence.
//
// A second return value indicates whether the sequence contained any values.
func Last[V any](seq iter.Seq[V]) (V, bool) {
	var v V
	var found bool

	for v = range seq {
		found = true
	}

	return v, found
}

// LastFunc returns the last value of a sequence that satisfies a predicate.
//
// A second return value indicates whether the sequence contained any value that satisfied the predicate.
func LastFunc[V any](seq iter.Seq[V], f func(V) bool) (V, bool) {
	var last V
	var found bool

	for v := range seq {
		if f(v) {
			last = v
			found = true
		}
	}

	return last, found
}

// Max returns the maximum value in a sequence.
//
// A second return value indicates whether the sequence contained any values.
func Max[V cmp.Ordered](seq iter.Seq[V]) (V, bool) {
	var maxVal V
	var found bool

	for v := range seq {
		if !found || v > maxVal {
			maxVal = v
		}

		found = true
	}

	return maxVal, found
}

// MaxBy returns the maximum value in a sequence using a function to select a key to use for comparisons.
//
// A second return value indicates whether the sequence contained any values.
func MaxBy[V any, K cmp.Ordered](seq iter.Seq[V], f func(V) K) (V, bool) {
	var maxK K
	var maxVal V
	var found bool

	for v := range seq {
		k := f(v)
		if !found || k > maxK {
			maxK = k
			maxVal = v
		}

		found = true
	}

	return maxVal, found
}

// MaxFunc returns the maximum value in a sequence using a comparison function.
//
// A second return value indicates whether the sequence contained any values.
func MaxFunc[V any](seq iter.Seq[V], f func(V, V) int) (V, bool) {
	var maxVal V
	var found bool

	for v := range seq {
		if !found || f(v, maxVal) > 0 {
			maxVal = v
		}

		found = true
	}

	return maxVal, found
}

// Min returns the minimum value in a sequence.
//
// A second return value indicates whether the sequence contained any values.
func Min[V cmp.Ordered](seq iter.Seq[V]) (V, bool) {
	var minVal V
	var found bool

	for v := range seq {
		if !found || v < minVal {
			minVal = v
		}

		found = true
	}

	return minVal, found
}

// MinBy returns the minimum value in a sequence using a function to select a key to use for comparisons.
func MinBy[V any, K cmp.Ordered](seq iter.Seq[V], f func(V) K) (V, bool) {
	var minK K
	var minVal V
	var found bool

	for v := range seq {
		k := f(v)
		if !found || k < minK {
			minK = k
			minVal = v
		}

		found = true
	}

	return minVal, found
}

// MinFunc returns the minimum value in a sequence using a comparison function.
//
// A second return value indicates whether the sequence contained any values.
func MinFunc[V any](seq iter.Seq[V], f func(V, V) int) (V, bool) {
	var minVal V
	var found bool

	for v := range seq {
		if !found || f(v, minVal) < 0 {
			minVal = v
		}

		found = true
	}

	return minVal, found
}

// OfType filters a sequence based on a type.
func OfType[V, VOut any](seq iter.Seq[V]) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for v := range seq {
			if out, ok := any(v).(VOut); ok {
				if !yield(out) {
					return
				}
			}
		}
	}
}

// Prepend adds values to the beginning of a sequence.
func Prepend[V any](seq iter.Seq[V], vals ...V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}

		for v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

// Range returns a sequence of numbers from start to end (inclusive)
// incremented by a given step size.
//
// If iteration of the sequence is stopped early and resumed,
// the sequence will start from the beginning.
func Range[V constraints.Integer | constraints.Float](start, end, step V) (iter.Seq[V], error) {
	switch {
	case step == 0:
		return Empty[V](), errors.New("seq.Range: step must be non-zero")

	case end > start && float64(step) < float64(0):
		return Empty[V](), errors.New("seq.Range: step must be positive for ascending ranges")

	case end < start && float64(step) > float64(0):
		return Empty[V](), errors.New("seq.Range: step must be negative for descending ranges")
	}

	if end < start {
		// Descending
		return func(yield func(V) bool) {
			for i := start; i >= end; i += step {
				if !yield(i) {
					return
				}
			}
		}, nil
	}

	// Ascending or equal
	return func(yield func(V) bool) {
		for i := start; i <= end; i += step {
			if !yield(i) {
				return
			}
		}
	}, nil
}

// Repeat returns a sequence that yields a given value n numer of times.
func Repeat[V any](val V, n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for range n {
			if !yield(val) {
				return
			}
		}
	}
}

// Select projects each value of a sequence into a new value.
func Select[V, VOut any](seq iter.Seq[V], f func(V) VOut) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for v := range seq {
			out := f(v)
			if !yield(out) {
				return
			}
		}
	}
}

// SelectKeys projects each value of a sequence into a key-value pair.
func SelectKeys[K, V any](seq iter.Seq[V], f func(V) K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for v := range seq {
			k := f(v)
			if !yield(k, v) {
				return
			}
		}
	}
}

// SelectMany projects each value of a sequence into a sequence and then flattens the resulting sequences
// into a single sequence.
func SelectMany[V, VOut any](seq iter.Seq[V], f func(V) iter.Seq[VOut]) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for v := range seq {
			for out := range f(v) {
				if !yield(out) {
					return
				}
			}
		}
	}
}

// Single returns the only value in a sequence.
//
// A second return value indicates whether the sequence contained exactly one value.
func Single[V any](seq iter.Seq[V]) (V, bool) {
	var first V
	var found bool

	for v := range seq {
		if found {
			var zero V
			return zero, false
		}

		first = v
		found = true
	}

	return first, found
}

// SingleFunc returns the only value in a sequence that satisfies a predicate.
//
// A second return value indicates whether the sequence contained exactly one value
// that satisfied the predicate.
func SingleFunc[V any](seq iter.Seq[V], f func(V) bool) (V, bool) {
	var first V
	var found bool

	for v := range seq {
		if f(v) {
			if found {
				var zero V
				return zero, false
			}

			first = v
			found = true
		}
	}

	return first, found
}

// Skip bypasses a given number of values in a sequence and returns the remaining values.
func Skip[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if i >= n {
				if !yield(v) {
					return
				}
			}

			i++
		}
	}
}

// SkipWhile bypasses values in a sequence as long as a condition is true and then
// returns the remaining values.
func SkipWhile[V any](seq iter.Seq[V], f func(int, V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		skip := true
		i := -1

		for v := range seq {
			i++

			if skip && f(i, v) {
				continue
			}

			if !yield(v) {
				return
			}

			skip = false
		}
	}
}

// Sum computes the sum of the values in a sequence.
func Sum[V constraints.Integer | constraints.Float](seq iter.Seq[V]) V {
	var sum V
	for v := range seq {
		sum += v
	}

	return sum
}

// Take returns a given number of values from the start of a sequence.
func Take[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if i >= n {
				return
			}

			if !yield(v) {
				return
			}

			i++
		}
	}
}

// TakeWhile returns values from a sequence as long as a given condition is true
// and then skips the remaining values.
func TakeWhile[V any](seq iter.Seq[V], f func(int, V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if !f(i, v) {
				return
			}

			if !yield(v) {
				return
			}

			i++
		}
	}
}

// Where filters a sequence based on a predicate.
func Where[V any](seq iter.Seq[V], f func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if f(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// ValueAt returns the value at a given index in a sequence.
//
// A second return value indicates whether the given index was within the bounds of the sequence.
// This panics if the given index is negative.
func ValueAt[V any](seq iter.Seq[V], index int) (V, bool) {
	if index < 0 {
		panic("seq.ValueAt: index must be non-negative")
	}

	var zero V
	i := 0

	for v := range seq {
		if i == index {
			return v, true
		}

		i++
	}

	return zero, false
}

// Yield returns a sequence of values.
//
// This is useful for creating a sequence from a slice or variadic arguments.
//
// Examples:
//
//	// yield each element of the slice
//	// yields ("a"), ("b"), ("c")
//	letters := []string{"a", "b", "c"}
//	vals := seq.Yield(letters...)
//
//	// yields (1), (2), (3)
//	vals := seq.Yield(1, 2, 3)
func Yield[V any](vals ...V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
	}
}

// YieldBackwards returns a sequence of values in reverse order.
//
// Example:
//
//	// yields (3), (2), (1)
//	vals := seq.YieldBackwards(1, 2, 3)
func YieldBackwards[V any](vals ...V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := len(vals) - 1; i >= 0; i-- {
			if !yield(vals[i]) {
				return
			}
		}
	}
}

// YieldChan returns a sequence of values from a channel.
func YieldChan[V any](ch <-chan V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range ch {
			if !yield(v) {
				return
			}
		}
	}
}

// ToSlice converts a sequence to a slice.
//
// This iterates over the entire sequence to collect all values.
// Use with caution for very large sequences as it loads all values into memory.
func ToSlice[V any](seq iter.Seq[V]) []V {
	if seq == nil {
		return make([]V, 0)
	}

	result := make([]V, 0)
	for v := range seq {
		result = append(result, v)
	}
	return result
}

// ToMap converts a sequence of key-value pairs to a map.
//
// This iterates over the entire sequence to collect all key-value pairs.
// Use with caution for very large sequences as it loads all values into memory.
// If the sequence contains duplicate keys, later values will overwrite earlier ones.
func ToMap[K comparable, V any](seq iter.Seq2[K, V]) map[K]V {
	if seq == nil {
		return make(map[K]V)
	}

	result := make(map[K]V)
	for k, v := range seq {
		result[k] = v
	}
	return result
}
