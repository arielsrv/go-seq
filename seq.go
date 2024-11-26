package seq

import (
	"cmp"
	"iter"

	"golang.org/x/exp/constraints"
)

func Aggregate[V, A any](seq iter.Seq[V], init A, f func(A, V) A) A {
	acc := init
	for v := range seq {
		acc = f(acc, v)
	}
	return acc
}

func All[V any](seq iter.Seq[V], f func(V) bool) bool {
	for v := range seq {
		if !f(v) {
			return false
		}
	}
	return true
}

func Any[V any](seq iter.Seq[V]) bool {
	for range seq {
		return true
	}
	return false
}

func AnyFunc[V any](seq iter.Seq[V], f func(V) bool) bool {
	for v := range seq {
		if f(v) {
			return true
		}
	}
	return false
}

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

func Chunk[V any](seq iter.Seq[V], size int) iter.Seq[[]V] {
	return func(yield func([]V) bool) {
		chunk := make([]V, 0, size)
		for v := range seq {
			chunk = append(chunk, v)
			if len(chunk) == size {
				if !yield(chunk) {
					return
				}
				chunk = make([]V, 0, size)
			}
		}
		if len(chunk) > 0 {
			if !yield(chunk) {
				return
			}
		}
	}
}

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

func Contains[V comparable](seq iter.Seq[V], val V) bool {
	for v := range seq {
		if v == val {
			return true
		}
	}
	return false
}

func Count[V any](seq iter.Seq[V]) int {
	count := 0
	for range seq {
		count++
	}
	return count
}

func CountFunc[V any](seq iter.Seq[V], f func(V) bool) int {
	count := 0
	for v := range seq {
		if f(v) {
			count++
		}
	}
	return count
}

func Distinct[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		seen := make(map[V]struct{})
		for v := range seq {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}

func DistinctFunc[V any, K comparable](seq iter.Seq[V], f func(V) K) iter.Seq[V] {
	return func(yield func(V) bool) {
		seen := make(map[K]struct{})
		for v := range seq {
			k := f(v)
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}

func ElementAt[V any](seq iter.Seq[V], index int) (V, bool) {
	i := 0
	for v := range seq {
		if i == index {
			return v, true
		}
		i++
	}

	var v V
	return v, false
}

func Empty[V any]() iter.Seq[V] {
	return func(yield func(V) bool) {}
}

func First[V any](seq iter.Seq[V]) (V, bool) {
	var v V
	for v = range seq {
		return v, true
	}

	return v, false
}

func Last[V any](seq iter.Seq[V]) (V, bool) {
	var v V
	found := false
	for v = range seq {
		found = true
	}

	return v, found
}

func Max[V cmp.Ordered](seq iter.Seq[V]) (V, bool) {
	var max V
	var found bool
	for v := range seq {
		if !found || v > max {
			max = v
			found = true
		}
	}
	return max, found
}

func Min[V cmp.Ordered](seq iter.Seq[V]) (V, bool) {
	var min V
	var found bool
	for v := range seq {
		if !found || v < min {
			min = v
			found = true
		}
	}
	return min, found
}

func NewSeq[V any](vals ...V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
	}
}

func None[V any](seq iter.Seq[V]) bool {
	for range seq {
		return false
	}
	return true
}

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

func Range[V constraints.Integer](start, end, step V) iter.Seq[V] {
	if step < 1 {
		panic("step must be positive")
	}

	return func(yield func(V) bool) {
		for i := start; i <= end; i += step {
			if !yield(i) {
				return
			}
		}
	}
}

func Repeat[V any](val V, count int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := 0; i < count; i++ {
			if !yield(val) {
				return
			}
		}
	}
}

func Select[V, VOut any](seq iter.Seq[V], f func(V) VOut) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

func SelectIndex[V any](seq iter.Seq[V]) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		i := 0
		for v := range seq {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}

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

func SelectSlices[V, VOut any](seq iter.Seq[V], f func(V) []VOut) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for v := range seq {
			for _, out := range f(v) {
				if !yield(out) {
					return
				}
			}
		}
	}
}

func Single[V any](seq iter.Seq[V]) (V, bool) {
	var first V
	var found bool

	for v := range seq {
		if found {
			return first, false
		}
		first = v
		found = true
	}

	return first, true
}

func Skip[V any](seq iter.Seq[V], num int) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if i >= num {
				if !yield(v) {
					return
				}
			}
			i++
		}
	}
}

func SkipWhile[V any](seq iter.Seq[V], f func(int, V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		skip := true
		i := 0
		for v := range seq {
			if skip {
				if f(i, v) {
					continue
				}
				skip = false
			}
			if !yield(v) {
				return
			}
			i++
		}
	}
}

func Sum[V constraints.Integer | constraints.Float](seq iter.Seq[V]) V {
	var sum V
	for v := range seq {
		sum += v
	}
	return sum
}

func Take[V any](seq iter.Seq[V], num int) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if i >= num {
				return
			}
			if !yield(v) {
				return
			}
			i++
		}
	}
}

func TakeWhile[V any](seq iter.Seq[V], f func(int, V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if f(i, v) {
				if !yield(v) {
					return
				}
			} else {
				return
			}
			i++
		}
	}
}

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

func Zip[V1 any, V2 any](seq1 iter.Seq[V1], seq2 iter.Seq[V2]) iter.Seq2[V1, V2] {
	return func(yield func(V1, V2) bool) {
		next2, stop := iter.Pull(seq2)
		defer stop()

		for v1 := range seq1 {
			v2, ok := next2()
			if !ok {
				return
			}
			if !yield(v1, v2) {
				return
			}
		}
	}
}
