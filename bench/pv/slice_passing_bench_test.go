package pv

import (
	"testing"
)

// processUtxoSliceByValue accepts an entire slice of Utxo structs by value.
// In Go, this only copies the 24-byte slice header, not the underlying data.
//go:noinline
func processUtxoSliceByValue(s []Utxo) {
	// Perform trivial work on the slice to prevent optimizations.
	var totalAmount int64
	for _, u := range s {
		totalAmount += int64(u.Amount)
	}
	sinkI64 = totalAmount
}

// processUtxoSliceByPointer accepts an entire slice of Utxo pointers.
// This also only copies the 24-byte slice header.
//go:noinline
func processUtxoSliceByPointer(s []*Utxo) {
	// Perform trivial work on the slice to prevent optimizations.
	var totalAmount int64
	for _, u := range s {
		totalAmount += int64(u.Amount)
	}
	sinkI64 = totalAmount
}

// BenchmarkSlicePassingCost measures the cost of passing an entire slice to a
// function. It demonstrates that passing a slice of values (`[]Utxo`) has the
// same negligible cost as passing a slice of pointers (`[]*Utxo`), because in
// both cases only the small slice header is copied.
//
// This benchmark confirms that the performance concerns are not about passing
// or returning entire collections, but about how the elements of those
// collections are constructed and used individually.
func BenchmarkSlicePassingCost(b *testing.B) {
	const (
		numUtxos   = 1 << 12 // 4096 UTXOs
		scriptSize = 34
	)

	// Pre-build a constant script and the slices to be passed.
	pkScript := make([]byte, scriptSize)
	for j := 0; j < scriptSize; j++ {
		pkScript[j] = byte(j)
	}
	utxoValues := buildUtxoValues(numUtxos, pkScript)
	utxoPointers := buildUtxoPointers(numUtxos, pkScript)

	b.Run("4k-Utxos/0-Pass-Value-Slice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			processUtxoSliceByValue(utxoValues)
		}
	})

	b.Run("4k-Utxos/1-Pass-Pointer-Slice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			processUtxoSliceByPointer(utxoPointers)
		}
	})
}
