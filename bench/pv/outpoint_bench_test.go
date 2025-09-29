package pv

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func makeOutPointValue(i int) wire.OutPoint {
	return wire.OutPoint{Hash: chainhash.Hash{byte(i % 251)}, Index: uint32(i)}
}

func makeOutPointPointer(i int) *wire.OutPoint {
	return &wire.OutPoint{Hash: chainhash.Hash{byte(i % 251)}, Index: uint32(i)}
}

func buildOutPointValues(n int) []wire.OutPoint {
	s := make([]wire.OutPoint, n)
	for i := 0; i < n; i++ {
		s[i] = makeOutPointValue(i)
	}
	return s
}

func buildOutPointPointers(n int) []*wire.OutPoint {
	s := make([]*wire.OutPoint, n)
	for i := 0; i < n; i++ {
		s[i] = makeOutPointPointer(i)
	}
	return s
}

// BenchmarkOutPoint_SliceBuild benchmarks building slices of OutPoint values vs pointers
func BenchmarkOutPoint_SliceBuild(b *testing.B) {
	datasets := generateOutPointDatasets(outpointBenchConfig{
		outpointGrowth: scaleGrowth(8, exponentialGrowth()),
		iterations:     8,
	})
	for _, d := range datasets {
		prefix := d.name()
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var s []wire.OutPoint
			for b.Loop() {
				s = buildOutPointValues(d.numOutPoints)
			}
			sinkInt = len(s)
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var s []*wire.OutPoint
			for b.Loop() {
				s = buildOutPointPointers(d.numOutPoints)
			}
			sinkInt = len(s)
		})
	}
}

// BenchmarkOutPoint_SliceIterate benchmarks iterating over slices of OutPoint values vs pointers
func BenchmarkOutPoint_SliceIterate(b *testing.B) {
	datasets := generateOutPointDatasets(outpointBenchConfig{
		outpointGrowth: scaleGrowth(8, exponentialGrowth()),
		iterations:     8,
	})
	for _, d := range datasets {
		vals := buildOutPointValues(d.numOutPoints)
		ptrs := buildOutPointPointers(d.numOutPoints)
		prefix := d.name()
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				for _, val := range vals {
					acc += int64(val.Index) + int64(val.Hash[0])
				}
			}
			sinkI64 = acc
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				for _, ptr := range ptrs {
					acc += int64(ptr.Index) + int64(ptr.Hash[0])
				}
			}
			sinkI64 = acc
		})
	}
}

// BenchmarkOutPoint_SliceBuildAndIterate benchmarks building and iterating over slices
// of OutPoint values vs pointers with repeated reads.
func BenchmarkOutPoint_SliceBuildAndIterate(b *testing.B) {
	d := outpointDataset{
		numOutPoints: 128,
		maxOutPoints: 128,
	}
	for i := range 10 {
		prefix := d.name() + fmt.Sprintf("NReads%d", i)
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				vals := buildOutPointValues(d.numOutPoints)
				for range i {
					for _, val := range vals {
						acc += int64(val.Index) + int64(val.Hash[0])
					}
				}
			}
			sinkI64 = acc
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				ptrs := buildOutPointPointers(d.numOutPoints)
				for range i {
					for _, ptr := range ptrs {
						acc += int64(ptr.Index) + int64(ptr.Hash[0])
					}
				}
			}
			sinkI64 = acc
		})
	}
}
