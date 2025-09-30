package pv

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
)

// Utxo provides a detailed overview of an unspent transaction output.
type Utxo struct {
	// OutPoint is the transaction output identifier.
	OutPoint wire.OutPoint

	// Amount is the value of the output.
	Amount btcutil.Amount

	// PkScript is the public key script for the output.
	PkScript []byte

	// Confirmations is the number of confirmations the output has.
	Confirmations int32

	// Spendable indicates whether the output is considered spendable.
	Spendable bool

	// Address is the address associated with the output.
	Address btcutil.Address

	// Account is the name of the account that owns the output.
	Account string

	// AddressType is the type of the address.
	AddressType waddrmgr.AddressType

	// Locked indicates whether the output is locked.
	Locked bool
}

// makeUtxoValue creates a Utxo by value, using a pre-made pkScript.
func makeUtxoValue(i int, pkScript []byte) Utxo {
	return Utxo{
		OutPoint:      wire.OutPoint{Hash: chainhash.Hash{byte(i % 251)}, Index: uint32(i)},
		Amount:        btcutil.Amount(1000 + i),
		PkScript:      pkScript,
		Confirmations: int32(i % 100),
		Spendable:     i%2 == 0,
		Address:       nil,
		Account:       "default",
		AddressType:   waddrmgr.WitnessPubKey,
		Locked:        false,
	}
}

// makeUtxoPointer creates a Utxo by pointer, using a pre-made pkScript.
func makeUtxoPointer(i int, pkScript []byte) *Utxo {
	return &Utxo{
		OutPoint:      wire.OutPoint{Hash: chainhash.Hash{byte(i % 251)}, Index: uint32(i)},
		Amount:        btcutil.Amount(1000 + i),
		PkScript:      pkScript,
		Confirmations: int32(i % 100),
		Spendable:     i%2 == 0,
		Address:       nil,
		Account:       "default",
		AddressType:   waddrmgr.WitnessPubKey,
		Locked:        false,
	}
}

// buildUtxoValues constructs a slice of Utxo values. The pkScript is shared.
func buildUtxoValues(n int, pkScript []byte) []Utxo {
	var s []Utxo
	for i := 0; i < n; i++ {
		s = append(s, makeUtxoValue(i, pkScript))
	}
	return s
}

// buildUtxoPointers constructs a slice of Utxo pointers. The pkScript is shared.
func buildUtxoPointers(n int, pkScript []byte) []*Utxo {
	var s []*Utxo
	for i := 0; i < n; i++ {
		s = append(s, makeUtxoPointer(i, pkScript))
	}
	return s
}

// BenchmarkUtxo_SliceBuild benchmarks building slices of Utxo values vs pointers
func BenchmarkUtxo_SliceBuild(b *testing.B) {
	datasets := generateUtxoDatasets(utxoBenchConfig{
		utxoGrowth:   scaleGrowth(8, exponentialGrowth()),
		scriptGrowth: linearGrowth(34),
		iterations:   8,
	})
	for _, d := range datasets {
		prefix := d.name()
		// Pre-build a constant script to be shared among all Utxos.
		// This isolates the benchmark to the cost of building the slice
		// of structs, not the cost of building the scripts themselves.
		pkScript := make([]byte, d.scriptSize)
		for j := 0; j < d.scriptSize; j++ {
			pkScript[j] = byte(j)
		}

		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var s []Utxo
			for b.Loop() {
				s = buildUtxoValues(d.numUtxos, pkScript)
			}
			sinkInt = len(s)
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var s []*Utxo
			for b.Loop() {
				s = buildUtxoPointers(d.numUtxos, pkScript)
			}
			sinkInt = len(s)
		})
	}
}

// BenchmarkUtxo_SliceIterate benchmarks iterating over slices of Utxo values vs pointers
func BenchmarkUtxo_SliceIterate(b *testing.B) {
	datasets := generateUtxoDatasets(utxoBenchConfig{
		utxoGrowth:   scaleGrowth(8, exponentialGrowth()),
		scriptGrowth: linearGrowth(34),
		iterations:   8,
	})
	for _, d := range datasets {
		pkScript := make([]byte, d.scriptSize)
		for j := 0; j < d.scriptSize; j++ {
			pkScript[j] = byte(j)
		}
		vals := buildUtxoValues(d.numUtxos, pkScript)
		ptrs := buildUtxoPointers(d.numUtxos, pkScript)
		prefix := d.name()
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				for _, val := range vals {
					acc += int64(val.Amount) + int64(len(val.PkScript)) + int64(val.Confirmations)
				}
			}
			sinkI64 = acc
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				for _, ptr := range ptrs {
					acc += int64(ptr.Amount) + int64(len(ptr.PkScript)) + int64(ptr.Confirmations)
				}
			}
			sinkI64 = acc
		})
	}
}

// BenchmarkUtxo_SliceBuildAndIterate benchmarks building and iterating over slices of Utxo values vs pointers
func BenchmarkUtxo_SliceBuildAndIterate(b *testing.B) {
	d := utxoDataset{
		numUtxos:   128,
		scriptSize: 64,
		maxUtxos:   128,
		maxScript:  64,
	}
	pkScript := make([]byte, d.scriptSize)
	for j := 0; j < d.scriptSize; j++ {
		pkScript[j] = byte(j)
	}

	for i := range 10 {
		prefix := d.name() + fmt.Sprintf("NReads%d", i)
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				vals := buildUtxoValues(d.numUtxos, pkScript)
				for range i {
					for _, val := range vals {
						acc += int64(val.Amount) + int64(len(val.PkScript)) + int64(val.Confirmations)
					}
				}
			}
			sinkI64 = acc
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				ptrs := buildUtxoPointers(d.numUtxos, pkScript)
				for range i {
					for _, ptr := range ptrs {
						acc += int64(ptr.Amount) + int64(len(ptr.PkScript)) + int64(ptr.Confirmations)
					}
				}
			}
			sinkI64 = acc
		})
	}
}