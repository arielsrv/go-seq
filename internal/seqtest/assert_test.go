package seqtest

import (
	"iter"
	"testing"
)

func Test_AssertEqual(t *testing.T) {
	AssertEqual(t, []int{1, 2, 3}, iter.Seq[int](func(yield func(int) bool) {
		yield(1)
		yield(2)
		yield(3)
	}))
}

func Test_AssertEqual2(t *testing.T) {
	expected := []KeyValuePair[string, int]{{"a", 1}, {"b", 2}}
	seq := iter.Seq2[string, int](func(yield func(string, int) bool) {
		yield("a", 1)
		yield("b", 2)
	})
	AssertEqual2(t, expected, seq)
}

func Test_AssertElementsMatch2(t *testing.T) {
	expected := []KeyValuePair[string, int]{{"a", 1}, {"b", 2}}
	seq := iter.Seq2[string, int](func(yield func(string, int) bool) {
		yield("b", 2)
		yield("a", 1)
	})
	AssertElementsMatch2(t, expected, seq)
}
