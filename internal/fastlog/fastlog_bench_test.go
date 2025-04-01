package fastlog

import (
	"math"
	"testing"
)

func BenchmarkFastLog2(b *testing.B) {
	for range b.N {
		for example := range examplesSeq() {
			Log2(example.X)
		}
	}
}

func BenchmarkMathLog2(b *testing.B) {
	for range b.N {
		for example := range examplesSeq() {
			math.Log2(float64(example.X))
		}
	}
}
