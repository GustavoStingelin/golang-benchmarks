package pv

import (
    "fmt"
    "testing"

    "github.com/btcsuite/btcd/wire"
)

func makeTxOutValue(i int, scriptSize int) wire.TxOut {
    if scriptSize < 0 {
        scriptSize = 0
    }
    return wire.TxOut{
        Value:    int64(1000 + i),
        PkScript: make([]byte, scriptSize),
    }
}

func makeTxOutPointer(i int, scriptSize int) *wire.TxOut {
    if scriptSize < 0 {
        scriptSize = 0
    }
    return &wire.TxOut{
        Value:    int64(1000 + i),
        PkScript: make([]byte, scriptSize),
    }
}

func buildTxOutValues(n, scriptSize int) []wire.TxOut {
    s := make([]wire.TxOut, n)
    for i := 0; i < n; i++ {
        s[i] = makeTxOutValue(i, scriptSize)
    }
    return s
}

func buildTxOutPointers(n, scriptSize int) []*wire.TxOut {
    s := make([]*wire.TxOut, n)
    for i := 0; i < n; i++ {
        s[i] = makeTxOutPointer(i, scriptSize)
    }
    return s
}

// BenchmarkTxOut_SliceBuild benchmarks building slices of TxOut values vs pointers
func BenchmarkTxOut_SliceBuild(b *testing.B) {
    datasets := generateTxOutDatasets(txoutBenchConfig{
        txoutGrowth:  scaleGrowth(8, exponentialGrowth()),
        scriptGrowth: linearGrowth(34),
        iterations:   8,
    })
    for _, d := range datasets {
        prefix := d.name()
        b.Run(prefix+"/0-Values", func(b *testing.B) {
            b.ReportAllocs()
            var s []wire.TxOut
            for b.Loop() {
                s = buildTxOutValues(d.numTxOuts, d.scriptSize)
            }
            sinkInt = len(s)
        })
        b.Run(prefix+"/1-Pointers", func(b *testing.B) {
            b.ReportAllocs()
            var s []*wire.TxOut
            for b.Loop() {
                s = buildTxOutPointers(d.numTxOuts, d.scriptSize)
            }
            sinkInt = len(s)
        })
    }
}

// BenchmarkTxOut_SliceIterate benchmarks iterating over slices of TxOut values vs pointers
func BenchmarkTxOut_SliceIterate(b *testing.B) {
    datasets := generateTxOutDatasets(txoutBenchConfig{
        txoutGrowth:  scaleGrowth(8, exponentialGrowth()),
        scriptGrowth: linearGrowth(34),
        iterations:   8,
    })
    for _, d := range datasets {
        vals := buildTxOutValues(d.numTxOuts, d.scriptSize)
        ptrs := buildTxOutPointers(d.numTxOuts, d.scriptSize)
        prefix := d.name()
        b.Run(prefix+"/0-Values", func(b *testing.B) {
            b.ReportAllocs()
            var acc int64
            for b.Loop() {
                for _, val := range vals {
                    acc += val.Value + int64(len(val.PkScript))
                }
            }
            sinkI64 = acc
        })
        b.Run(prefix+"/1-Pointers", func(b *testing.B) {
            b.ReportAllocs()
            var acc int64
            for b.Loop() {
                for _, ptr := range ptrs {
                    acc += ptr.Value + int64(len(ptr.PkScript))
                }
            }
            sinkI64 = acc
        })
    }
}

// BenchmarkTxOut_SliceBuildAndIterate benchmarks building and iterating over slices
// of TxOut values vs pointers with repeated reads.
func BenchmarkTxOut_SliceBuildAndIterate(b *testing.B) {
    d := txoutDataset{
        numTxOuts:  128,
        scriptSize: 64,
        maxTxOuts:  128,
        maxScript:  64,
    }
    for i := range 10 {
        prefix := d.name() + fmt.Sprintf("NReads%d", i)
        b.Run(prefix+"/0-Values", func(b *testing.B) {
            b.ReportAllocs()
            var acc int64
            for b.Loop() {
                vals := buildTxOutValues(d.numTxOuts, d.scriptSize)
                for range i {
                    for _, val := range vals {
                        acc += val.Value + int64(len(val.PkScript))
                    }
                }
            }
            sinkI64 = acc
        })
        b.Run(prefix+"/1-Pointers", func(b *testing.B) {
            b.ReportAllocs()
            var acc int64
            for b.Loop() {
                ptrs := buildTxOutPointers(d.numTxOuts, d.scriptSize)
                for range i {
                    for _, ptr := range ptrs {
                        acc += ptr.Value + int64(len(ptr.PkScript))
                    }
                }
            }
            sinkI64 = acc
        })
    }
}

