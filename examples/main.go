package main

import (
	"fmt"
	"iter"
)

func main() {
	Range_Examples()
}

func PrintAll[V any](seq iter.Seq[V]) {
	for v := range seq {
		fmt.Println(v)
	}
}

func PrintRunes(seq iter.Seq[rune]) {
	for r := range seq {
		fmt.Println(string(r))
	}
}

func PrintAll2[K, V any](seq iter.Seq2[K, V]) {
	for k, v := range seq {
		fmt.Println(k, v)
	}
}
