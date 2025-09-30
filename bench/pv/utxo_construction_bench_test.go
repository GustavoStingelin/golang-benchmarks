package pv

import (
	"testing"
)

// BenchmarkUtxoConstruction isolates and measures the performance of
// constructing a slice of Utxo structs, comparing the value vs. pointer approach.
//
// This benchmark demonstrates the "cost of construction":
//   - Values ([]Utxo): Building the slice requires copying the entire Utxo
//     struct for each element that is appended. As the slice grows, this
//     leads to significant memory being copied, which is reflected in the
//     high B/op (bytes per operation) metric.
//   - Pointers ([]*Utxo): Building the slice only requires copying an 8-byte
//     pointer for each element. This results in far less memory being copied,
//     but it requires a separate heap allocation for each Utxo, leading to a
//     higher allocs/op metric.
//
// For a large and complex struct like Utxo, the high memory copy overhead
// for the value slice makes it a less performant choice during the build phase.
func BenchmarkUtxoConstruction(b *testing.B) {
	const (
		numUtxos   = 1 << 12 // 4096 UTXOs
		scriptSize = 34      // A typical P2WKH script size
	)

	// Pre-build a constant script to be shared among all Utxos.
	// This isolates the benchmark to the cost of building the slice
	// of structs, not the cost of building the scripts themselves.
	pkScript := make([]byte, scriptSize)
	for j := 0; j < scriptSize; j++ {
		pkScript[j] = byte(j)
	}

	b.Run("4k-Utxos/0-Values-Construction", func(b *testing.B) {
		b.ReportAllocs()
		var s []Utxo
		for i := 0; i < b.N; i++ {
			// In each iteration, we build the slice from scratch.
			s = buildUtxoValues(numUtxos, pkScript)
		}
		// Prevent the compiler from optimizing away the slice.
		sinkInt = len(s)
	})

	b.Run("4k-Utxos/1-Pointers-Construction", func(b *testing.B) {
		b.ReportAllocs()
		var s []*Utxo
		for i := 0; i < b.N; i++ {
			s = buildUtxoPointers(numUtxos, pkScript)
		}
		sinkInt = len(s)
	})
}