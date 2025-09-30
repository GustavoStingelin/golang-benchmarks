package pv

import (
	"testing"
)

// processUtxoValue accepts a Utxo by value, forcing a copy of the
// entire struct onto the function's stack frame for each call.
//go:noinline
func processUtxoValue(u Utxo) {
	// Perform trivial work on the struct to prevent the compiler from
	// optimizing the function call away.
	sinkI64 += int64(u.Amount)
	if u.Spendable {
		sinkInt++
	}
}

// processUtxoPointer accepts a Utxo by pointer, copying only the 8-byte
// pointer for each call.
//go:noinline
func processUtxoPointer(u *Utxo) {
	// Perform trivial work on the struct.
	sinkI64 += int64(u.Amount)
	if u.Spendable {
		sinkInt++
	}
}

// BenchmarkUsageCost measures the performance impact of passing structs to
// functions by value vs. by pointer. It simulates iterating over a slice of
// UTXOs and calling a processing function for each one.
//
// The results will demonstrate the "cost of use":
//  - Values: Will be significantly slower because the entire Utxo struct
//    (approx. 104 bytes + slice header) is copied for every function call.
//  - Pointers: Will be much faster because only an 8-byte pointer is copied
//    for each call.
//
// This benchmark highlights why returning a slice of pointers can be more
// performant in real-world applications, even if raw iteration over a value
// slice is faster. The cost of actually *using* the elements often dominates.
func BenchmarkUsageCost(b *testing.B) {
	const (
		numUtxos   = 1 << 12 // 4096 UTXOs
		scriptSize = 34      // A typical P2WKH script size
	)

	// Pre-build the slices once, outside the timed loop.
	pkScript := make([]byte, scriptSize)
	for j := 0; j < scriptSize; j++ {
		pkScript[j] = byte(j)
	}
	utxoValues := buildUtxoValues(numUtxos, pkScript)
	utxoPointers := buildUtxoPointers(numUtxos, pkScript)

	b.Run("4k-Utxos/0-Values", func(b *testing.B) {
		b.ReportAllocs()
		// The b.N loop is inside the for...range loop to ensure we are
		// measuring the cost of many function calls.
		for i := 0; i < b.N; i++ {
			for _, u := range utxoValues {
				processUtxoValue(u)
			}
		}
	})

	b.Run("4k-Utxos/1-Pointers", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, u := range utxoPointers {
				processUtxoPointer(u)
			}
		}
	})
}