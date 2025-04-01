package fastlog

import "math"

// Adapted from lc0/src/utils/fastmath.h

// Best known constant for K tuned using tune_constants.go.
// Tested using suite.go.
// Suite MSE for Log2: 2.9706729847024478e-05.
const (
	K  = 0.3462012
	c0 = 1 + K
)

// fastLog2 implements fast approximate Log2 using the MSE constants.
// Does no range checking.
func fastLog2(x, c0, c1 float32) float32 {
	tmp := math.Float32bits(x)
	expb := uint64(tmp) >> 23
	tmp = (tmp & 0x7fffff) | (0x7f << 23)
	out := math.Float32frombits(tmp) - 1
	return out*(c0+c1*out) - 127 + float32(expb)
}

// Log2 implements fast approximate Log2. Does no range checking.
func Log2(x float32) float32 {
	tmp := math.Float32bits(x)
	expb := uint64(tmp) >> 23
	tmp = (tmp & 0x7fffff) | (0x7f << 23)
	out := math.Float32frombits(tmp) - 1
	return out*(c0-K*out) - 127 + float32(expb)
}

// Log implements fast approximate Ln(x). Does no range checking.
func Log(x float32) float32 { return 0.6931471805599453 * Log2(x) }
