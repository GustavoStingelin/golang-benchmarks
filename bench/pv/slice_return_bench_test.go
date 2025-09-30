package pv

import (
	"fmt"
	"testing"
)

// returnUtxoValues and returnUtxoPointers are marked noinline to ensure the
// benchmark measures the cost of returning the slice header and does not get
// optimized away. Returning a slice only copies the 3-word slice header, not
// the backing array or any element data.

//go:noinline
func returnUtxoValues(s []Utxo) []Utxo {
	s[0].Amount++
	return s
}

//go:noinline
func returnUtxoPointers(s []*Utxo) []*Utxo {
	s[0].Amount++
	return s
}

// BenchmarkUtxo_ReturnOnly demonstrates that returning a slice (values or
// pointers) does not copy the underlying backing array. We pre-build a large
// slice (several MB across element data due to PkScript allocations) once
// outside the timing loop, then repeatedly return it from a noinline function.
// Allocations and bytes/op should be ~0 for both cases, disproving the claim
// that returning []T copies megabytes of data.
func BenchmarkUtxo_ReturnOnly(b *testing.B) {
	// Choose a size large enough to represent multiple megabytes when using
	// []Utxo values with non-trivial PkScript sizes, while keeping memory
	// usage reasonable for CI.
	const (
		n          = 1 << 15 // 32768 elements
		scriptSize = 256     // bytes per PkScript backing array
	)

	// Build once outside the timed loop.
	vals := buildUtxoValues(n, scriptSize)
	ptrs := buildUtxoPointers(n, scriptSize)

	name := fmt.Sprintf("%d-Utxo-%dScript-ReturnOnly", n, scriptSize)

	b.Run(name+"/0-Values", func(b *testing.B) {
		b.ReportAllocs()
		var s []Utxo
		for b.Loop() {
			s = returnUtxoValues(vals)
		}
		sinkInt = len(s)
	})

	b.Run(name+"/1-Pointers", func(b *testing.B) {
		b.ReportAllocs()
		var s []*Utxo
		for b.Loop() {
			s = returnUtxoPointers(ptrs)
		}
		sinkInt = len(s)
	})
}
