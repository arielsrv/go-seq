package main

import (
	"github.com/arielsrv/go-seq"
)

func RangeExamples() {
	numbers, _ := seq.Range(1, 10, 1)
	PrintAll(numbers)

	numbersBackwards, _ := seq.Range(10, 1, 1)
	PrintAll(numbersBackwards)

	floats, _ := seq.Range(0.0, 5.0, 0.5)
	PrintAll(floats)

	alphabet, _ := seq.Range('a', 'z', 1)
	PrintRunes(alphabet)

	alphaBackwards, _ := seq.Range('z', 'a', 1)
	PrintRunes(alphaBackwards)
}
