package pv

import (
	"testing"
)

// LargeStruct is a struct with a sizable data payload to make
// copy costs significant during slice construction.
type LargeStruct struct {
	ID   int64
	Data [256]byte // 256 bytes of data
}

// buildLargeStructValues constructs a slice of LargeStruct by value.
// The cost of this function is dominated by the copying of structs
// during the append() calls.
func buildLargeStructValues(n int) []LargeStruct {
	// Pre-allocating capacity with make([]T, 0, n) is a best practice,
	// but we use a zero-value slice here to better simulate a common
	// pattern where the final size is unknown and the slice must grow.
	var s []LargeStruct
	for i := 0; i < n; i++ {
		val := LargeStruct{
			ID: int64(i),
		}
		// This append() copies the entire 'val' struct (264 bytes)
		// into the slice's backing array.
		s = append(s, val)
	}
	return s
}

// buildLargeStructPointers constructs a slice of *LargeStruct by pointer.
// The cost of this function is dominated by the heap allocations for each
// struct, but the append() calls remain cheap.
func buildLargeStructPointers(n int) []*LargeStruct {
	var s []*LargeStruct
	for i := 0; i < n; i++ {
		val := &LargeStruct{
			ID: int64(i),
		}
		// This append() copies only the 8-byte pointer into the slice's
		// backing array.
		s = append(s, val)
	}
	return s
}

// BenchmarkConstructionCost provides a clear demonstration of the performance
// difference between constructing a slice of values vs. a slice of pointers.
//
// The results will show:
//   - Values: High ns/op (time), very high B/op (bytes copied), low allocs/op.
//   - Pointers: Low ns/op, low B/op, but higher allocs/op (one per element).
//
// For large structs, the massive byte-copying cost for the value slice
// makes it significantly slower overall.
func BenchmarkConstructionCost(b *testing.B) {
	const numElements = 10000

	b.Run("10k-Elements/0-Values", func(b *testing.B) {
		b.ReportAllocs()
		var s []LargeStruct
		for b.Loop() {
			s = buildLargeStructValues(numElements)
		}
		sinkInt = len(s)
	})

	b.Run("10k-Elements/1-Pointers", func(b *testing.B) {
		b.ReportAllocs()
		var s []*LargeStruct
		for b.Loop() {
			s = buildLargeStructPointers(numElements)
		}
		sinkInt = len(s)
	})
}
