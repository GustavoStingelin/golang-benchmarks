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

func makeUtxoValue(i int, scriptSize int) Utxo {
	if scriptSize < 0 {
		scriptSize = 0
	}
	return Utxo{
		OutPoint:      wire.OutPoint{Hash: chainhash.Hash{byte(i % 251)}, Index: uint32(i)},
		Amount:        btcutil.Amount(1000 + i),
		PkScript:      make([]byte, scriptSize),
		Confirmations: int32(i % 100),
		Spendable:     i%2 == 0,
		Address:       nil,
		Account:       "default",
		AddressType:   waddrmgr.WitnessPubKey,
		Locked:        false,
	}
}

func makeUtxoPointer(i int, scriptSize int) *Utxo {
	if scriptSize < 0 {
		scriptSize = 0
	}
	return &Utxo{
		OutPoint:      wire.OutPoint{Hash: chainhash.Hash{byte(i % 251)}, Index: uint32(i)},
		Amount:        btcutil.Amount(1000 + i),
		PkScript:      make([]byte, scriptSize),
		Confirmations: int32(i % 100),
		Spendable:     i%2 == 0,
		Address:       nil,
		Account:       "default",
		AddressType:   waddrmgr.WitnessPubKey,
		Locked:        false,
	}
}

func buildUtxoValues(n, scriptSize int) []Utxo {
	s := make([]Utxo, n)
	for i := 0; i < n; i++ {
		s[i] = makeUtxoValue(i, scriptSize)
	}
	return s
}

func buildUtxoPointers(n, scriptSize int) []*Utxo {
	s := make([]*Utxo, n)
	for i := 0; i < n; i++ {
		s[i] = makeUtxoPointer(i, scriptSize)
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
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var s []Utxo
			for b.Loop() {
				s = buildUtxoValues(d.numUtxos, d.scriptSize)
			}
			sinkInt = len(s)
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var s []*Utxo
			for b.Loop() {
				s = buildUtxoPointers(d.numUtxos, d.scriptSize)
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
		vals := buildUtxoValues(d.numUtxos, d.scriptSize)
		ptrs := buildUtxoPointers(d.numUtxos, d.scriptSize)
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
	for i := range 10 {
		prefix := d.name() + fmt.Sprintf("NReads%d", i)
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				vals := buildUtxoValues(d.numUtxos, d.scriptSize)
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
				ptrs := buildUtxoPointers(d.numUtxos, d.scriptSize)
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
