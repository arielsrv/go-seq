package main

import (
	"testing"

	"iter"

	"github.com/arielsrv/go-seq"
	"github.com/arielsrv/go-seq/internal/seqtest"
)

func Test_Append_EdgeCases(t *testing.T) {
	t.Run("append to nil sequence", func(t *testing.T) {
		var nilSeq iter.Seq[int]
		got := seq.Append(nilSeq, 1, 2)
		seqtest.AssertEqual(t, []int{1, 2}, got)
	})
}

func Test_Prepend_EdgeCases(t *testing.T) {
	t.Run("prepend to nil sequence", func(t *testing.T) {
		var nilSeq iter.Seq[int]
		got := seq.Prepend(nilSeq, 1, 2)
		seqtest.AssertEqual(t, []int{1, 2}, got)
	})
}

func Test_Repeat_EdgeCases(t *testing.T) {
	t.Run("repeat zero times", func(t *testing.T) {
		got := seq.Repeat(42, 0)
		seqtest.AssertEqual(t, nil, got)
	})
	t.Run("repeat negative times", func(t *testing.T) {
		got := seq.Repeat(42, -5)
		seqtest.AssertEqual(t, nil, got)
	})
}

func Test_YieldBackwards_EdgeCases(t *testing.T) {
	t.Run("yield backwards nil", func(t *testing.T) {
		var nilVals []int
		got := seq.YieldBackwards(nilVals...)
		seqtest.AssertEqual(t, nil, got)
	})
}

func Test_YieldChan_EdgeCases(t *testing.T) {
	t.Run("yield chan nil", func(t *testing.T) {
		var ch chan int
		got := seq.YieldChan(ch)
		seqtest.AssertEqual(t, nil, got)
	})
}
