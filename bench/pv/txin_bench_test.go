package pv

import (
    "fmt"
    "testing"

    "github.com/btcsuite/btcd/chaincfg/chainhash"
    "github.com/btcsuite/btcd/wire"
)

func makeTxInValue(i int, scriptSize int) wire.TxIn {
    if scriptSize < 0 {
        scriptSize = 0
    }
    return wire.TxIn{
        PreviousOutPoint: wire.OutPoint{Hash: chainhash.Hash{byte(i % 251)}, Index: uint32(i)},
        SignatureScript:  make([]byte, scriptSize),
        Witness:          nil,
        Sequence:         uint32(100000 + i),
    }
}

func makeTxInPointer(i int, scriptSize int) *wire.TxIn {
    if scriptSize < 0 {
        scriptSize = 0
    }
    return &wire.TxIn{
        PreviousOutPoint: wire.OutPoint{Hash: chainhash.Hash{byte(i % 251)}, Index: uint32(i)},
        SignatureScript:  make([]byte, scriptSize),
        Witness:          nil,
        Sequence:         uint32(100000 + i),
    }
}

func buildTxInValues(n, scriptSize int) []wire.TxIn {
    s := make([]wire.TxIn, n)
    for i := 0; i < n; i++ {
        s[i] = makeTxInValue(i, scriptSize)
    }
    return s
}

func buildTxInPointers(n, scriptSize int) []*wire.TxIn {
    s := make([]*wire.TxIn, n)
    for i := 0; i < n; i++ {
        s[i] = makeTxInPointer(i, scriptSize)
    }
    return s
}

// BenchmarkTxIn_SliceBuild benchmarks building slices of TxIn values vs pointers
func BenchmarkTxIn_SliceBuild(b *testing.B) {
    datasets := generateTxInDatasets(txinBenchConfig{
        txinGrowth:   scaleGrowth(8, exponentialGrowth()),
        scriptGrowth: linearGrowth(34),
        iterations:   8,
    })
    for _, d := range datasets {
        prefix := d.name()
        b.Run(prefix+"/0-Values", func(b *testing.B) {
            b.ReportAllocs()
            var s []wire.TxIn
            for b.Loop() {
                s = buildTxInValues(d.numTxIns, d.scriptSize)
            }
            sinkInt = len(s)
        })
        b.Run(prefix+"/1-Pointers", func(b *testing.B) {
            b.ReportAllocs()
            var s []*wire.TxIn
            for b.Loop() {
                s = buildTxInPointers(d.numTxIns, d.scriptSize)
            }
            sinkInt = len(s)
        })
    }
}

// BenchmarkTxIn_SliceIterate benchmarks iterating over slices of TxIn values vs pointers
func BenchmarkTxIn_SliceIterate(b *testing.B) {
    datasets := generateTxInDatasets(txinBenchConfig{
        txinGrowth:   scaleGrowth(8, exponentialGrowth()),
        scriptGrowth: linearGrowth(34),
        iterations:   8,
    })
    for _, d := range datasets {
        vals := buildTxInValues(d.numTxIns, d.scriptSize)
        ptrs := buildTxInPointers(d.numTxIns, d.scriptSize)
        prefix := d.name()
        b.Run(prefix+"/0-Values", func(b *testing.B) {
            b.ReportAllocs()
            var acc int64
            for b.Loop() {
                for _, val := range vals {
                    acc += int64(len(val.SignatureScript))
                    acc += int64(val.Sequence)
                    acc += int64(val.PreviousOutPoint.Index)
                    acc += int64(val.PreviousOutPoint.Hash[0])
                }
            }
            sinkI64 = acc
        })
        b.Run(prefix+"/1-Pointers", func(b *testing.B) {
            b.ReportAllocs()
            var acc int64
            for b.Loop() {
                for _, ptr := range ptrs {
                    acc += int64(len(ptr.SignatureScript))
                    acc += int64(ptr.Sequence)
                    acc += int64(ptr.PreviousOutPoint.Index)
                    acc += int64(ptr.PreviousOutPoint.Hash[0])
                }
            }
            sinkI64 = acc
        })
    }
}

// BenchmarkTxIn_SliceBuildAndIterate benchmarks building and iterating over slices
// of TxIn values vs pointers with repeated reads.
func BenchmarkTxIn_SliceBuildAndIterate(b *testing.B) {
    d := txinDataset{
        numTxIns:   128,
        scriptSize: 64,
        maxTxIns:   128,
        maxScript:  64,
    }
    for i := range 10 {
        prefix := d.name() + fmt.Sprintf("NReads%d", i)
        b.Run(prefix+"/0-Values", func(b *testing.B) {
            b.ReportAllocs()
            var acc int64
            for b.Loop() {
                vals := buildTxInValues(d.numTxIns, d.scriptSize)
                for range i {
                    for _, val := range vals {
                        acc += int64(len(val.SignatureScript))
                        acc += int64(val.Sequence)
                        acc += int64(val.PreviousOutPoint.Index)
                        acc += int64(val.PreviousOutPoint.Hash[0])
                    }
                }
            }
            sinkI64 = acc
        })
        b.Run(prefix+"/1-Pointers", func(b *testing.B) {
            b.ReportAllocs()
            var acc int64
            for b.Loop() {
                ptrs := buildTxInPointers(d.numTxIns, d.scriptSize)
                for range i {
                    for _, ptr := range ptrs {
                        acc += int64(len(ptr.SignatureScript))
                        acc += int64(ptr.Sequence)
                        acc += int64(ptr.PreviousOutPoint.Index)
                        acc += int64(ptr.PreviousOutPoint.Hash[0])
                    }
                }
            }
            sinkI64 = acc
        })
    }
}

