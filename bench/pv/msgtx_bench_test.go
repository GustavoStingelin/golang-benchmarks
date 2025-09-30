package pv

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func makeMsgTxValue(i, nIn, nOut, scriptSize int) wire.MsgTx {
	if scriptSize < 0 {
		scriptSize = 0
	}
	tx := wire.MsgTx{
		Version:  2,
		TxIn:     make([]*wire.TxIn, nIn),
		TxOut:    make([]*wire.TxOut, nOut),
		LockTime: 0,
	}
	for j := 0; j < nIn; j++ {
		tx.TxIn[j] = &wire.TxIn{
			PreviousOutPoint: wire.OutPoint{Hash: chainhash.Hash{byte((i + j) % 251)}, Index: uint32(i + j)},
			SignatureScript:  make([]byte, scriptSize),
			Witness:          nil,
			Sequence:         uint32(100000 + i + j),
		}
	}
	for k := 0; k < nOut; k++ {
		tx.TxOut[k] = &wire.TxOut{
			Value:    int64(1000 + i + k),
			PkScript: make([]byte, scriptSize),
		}
	}
	return tx
}

func makeMsgTxPointer(i, nIn, nOut, scriptSize int) *wire.MsgTx {
	if scriptSize < 0 {
		scriptSize = 0
	}
	tx := wire.MsgTx{
		Version:  2,
		TxIn:     make([]*wire.TxIn, nIn),
		TxOut:    make([]*wire.TxOut, nOut),
		LockTime: 0,
	}
	for j := 0; j < nIn; j++ {
		tx.TxIn[j] = &wire.TxIn{
			PreviousOutPoint: wire.OutPoint{Hash: chainhash.Hash{byte((i + j) % 251)}, Index: uint32(i + j)},
			SignatureScript:  make([]byte, scriptSize),
			Witness:          nil,
			Sequence:         uint32(100000 + i + j),
		}
	}
	for k := 0; k < nOut; k++ {
		tx.TxOut[k] = &wire.TxOut{
			Value:    int64(1000 + i + k),
			PkScript: make([]byte, scriptSize),
		}
	}
	return &tx
}

func buildMsgTxValues(n, nIn, nOut, scriptSize int) []wire.MsgTx {
	s := make([]wire.MsgTx, n)
	for i := 0; i < n; i++ {
		s[i] = makeMsgTxValue(i, nIn, nOut, scriptSize)
	}
	return s
}

func buildMsgTxPointers(n, nIn, nOut, scriptSize int) []*wire.MsgTx {
	s := make([]*wire.MsgTx, n)
	for i := 0; i < n; i++ {
		s[i] = makeMsgTxPointer(i, nIn, nOut, scriptSize)
	}
	return s
}

// BenchmarkMsgTx_SliceBuild benchmarks building slices of MsgTx values vs pointers
func BenchmarkMsgTx_SliceBuild(b *testing.B) {
	datasets := generateMsgTxDatasets(msgtxBenchConfig{
		txGrowth:     scaleGrowth(4, exponentialGrowth()),
		scriptGrowth: linearGrowth(34),
		iterations:   8,
		nInputs:      2,
		nOutputs:     2,
	})
	for _, d := range datasets {
		prefix := d.name()
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var s []wire.MsgTx
			for b.Loop() {
				s = buildMsgTxValues(d.numTxs, d.nInputs, d.nOutputs, d.scriptSize)
			}
			sinkInt = len(s)
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var s []*wire.MsgTx
			for b.Loop() {
				s = buildMsgTxPointers(d.numTxs, d.nInputs, d.nOutputs, d.scriptSize)
			}
			sinkInt = len(s)
		})
	}
}

// BenchmarkMsgTx_SliceIterate benchmarks iterating over slices of MsgTx values vs pointers
func BenchmarkMsgTx_SliceIterate(b *testing.B) {
	datasets := generateMsgTxDatasets(msgtxBenchConfig{
		txGrowth:     scaleGrowth(4, exponentialGrowth()),
		scriptGrowth: linearGrowth(34),
		iterations:   8,
		nInputs:      2,
		nOutputs:     2,
	})
	for _, d := range datasets {
		vals := buildMsgTxValues(d.numTxs, d.nInputs, d.nOutputs, d.scriptSize)
		ptrs := buildMsgTxPointers(d.numTxs, d.nInputs, d.nOutputs, d.scriptSize)
		prefix := d.name()
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				for _, tx := range vals {
					// Sum across inputs and outputs to exercise nested fields.
					for _, ti := range tx.TxIn {
						acc += int64(len(ti.SignatureScript))
						acc += int64(ti.Sequence)
						acc += int64(ti.PreviousOutPoint.Index)
						acc += int64(ti.PreviousOutPoint.Hash[0])
					}
					for _, to := range tx.TxOut {
						acc += to.Value + int64(len(to.PkScript))
					}
				}
			}
			sinkI64 = acc
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				for _, tx := range ptrs {
					for _, ti := range tx.TxIn {
						acc += int64(len(ti.SignatureScript))
						acc += int64(ti.Sequence)
						acc += int64(ti.PreviousOutPoint.Index)
						acc += int64(ti.PreviousOutPoint.Hash[0])
					}
					for _, to := range tx.TxOut {
						acc += to.Value + int64(len(to.PkScript))
					}
				}
			}
			sinkI64 = acc
		})
	}
}

// BenchmarkMsgTx_SliceBuildAndIterate benchmarks building and iterating over slices
// of MsgTx values vs pointers with repeated reads.
func BenchmarkMsgTx_SliceBuildAndIterate(b *testing.B) {
	d := msgtxDataset{
		numTxs:     256,
		scriptSize: 64,
		nInputs:    2,
		nOutputs:   2,
		maxTxs:     256,
		maxScript:  64,
	}
	for i := range 10 {
		prefix := d.name() + fmt.Sprintf("NReads%d", i)
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				vals := buildMsgTxValues(d.numTxs, d.nInputs, d.nOutputs, d.scriptSize)
				for range i {
					for _, tx := range vals {
						for _, ti := range tx.TxIn {
							acc += int64(len(ti.SignatureScript))
							acc += int64(ti.Sequence)
							acc += int64(ti.PreviousOutPoint.Index)
							acc += int64(ti.PreviousOutPoint.Hash[0])
						}
						for _, to := range tx.TxOut {
							acc += to.Value + int64(len(to.PkScript))
						}
					}
				}
			}
			sinkI64 = acc
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				ptrs := buildMsgTxPointers(d.numTxs, d.nInputs, d.nOutputs, d.scriptSize)
				for range i {
					for _, tx := range ptrs {
						for _, ti := range tx.TxIn {
							acc += int64(len(ti.SignatureScript))
							acc += int64(ti.Sequence)
							acc += int64(ti.PreviousOutPoint.Index)
							acc += int64(ti.PreviousOutPoint.Hash[0])
						}
						for _, to := range tx.TxOut {
							acc += to.Value + int64(len(to.PkScript))
						}
					}
				}
			}
			sinkI64 = acc
		})
	}
}
