package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("var suite = []example{")
	// Add numbers 1..256 to the suite.
	for i := range 256 {
		outputLine(float32(i + 1))
	}
	// Add numbers [300 - 1000) to the suite.
	for x := float32(300); x < 1000; x *= 1.02 {
		outputLine(x)
	}
	// Add medium numbers to the suite.
	for x := float32(1000); x < 100_000; x *= 1.01 {
		outputLine(x)
	}
	// Add large numbers to the suite.
	for x := float32(100_000); x < 10_000_000; x *= 1.1 {
		outputLine(x)
	}
	// Add 1e7 to the suite.
	outputLine(10_000_000)
	fmt.Println("}")
}

func outputLine(x float32) {
	x = float32(int64(x)) // Truncate decimals.
	log2 := float32(math.Log2(float64(x)))
	log := float32(math.Log(float64(x)))
	// fmt.Printf("    makeExample(%d, %d, %d),\n", math.Float32bits(x), math.Float32bits(log2), math.Float32bits(log))
	fmt.Printf("    {%v, %v, %v},\n", x, log2, log)
}
