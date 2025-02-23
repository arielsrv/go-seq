package main

import (
	"github.com/sectrean/go-seq"
)

func Range_Examples() {
	numbers := seq.Range(1, 10, 1)
	PrintAll(numbers)

	numbersBackwards := seq.Range(10, 1, 1)
	PrintAll(numbersBackwards)

	floats := seq.Range(0.0, 5.0, 0.5)
	PrintAll(floats)

	alphabet := seq.Range('a', 'z', 1)
	PrintRunes(alphabet)

	alphaBackwards := seq.Range('z', 'a', 1)
	PrintRunes(alphaBackwards)
}
