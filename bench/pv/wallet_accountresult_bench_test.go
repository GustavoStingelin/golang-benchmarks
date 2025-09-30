package pv

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/wallet"
)

func makeAccountResultValue(i int) wallet.AccountResult {
	return wallet.AccountResult{
		AccountProperties: waddrmgr.AccountProperties{
			AccountNumber:        uint32(i),
			AccountName:          fmt.Sprintf("acct-%d", i),
			ExternalKeyCount:     uint32(i % 100),
			InternalKeyCount:     uint32((i + 7) % 100),
			ImportedKeyCount:     uint32((i + 13) % 100),
			AccountPubKey:        nil,
			MasterKeyFingerprint: uint32(i * 3),
			KeyScope:             waddrmgr.KeyScope{},
			IsWatchOnly:          i%2 == 0,
			AddrSchema:           nil,
		},
		TotalBalance: btcutil.Amount(1000 + i),
	}
}

func makeAccountResultPointer(i int) *wallet.AccountResult {
	return &wallet.AccountResult{
		AccountProperties: waddrmgr.AccountProperties{
			AccountNumber:        uint32(i),
			AccountName:          fmt.Sprintf("acct-%d", i),
			ExternalKeyCount:     uint32(i % 100),
			InternalKeyCount:     uint32((i + 7) % 100),
			ImportedKeyCount:     uint32((i + 13) % 100),
			AccountPubKey:        nil,
			MasterKeyFingerprint: uint32(i * 3),
			KeyScope:             waddrmgr.KeyScope{},
			IsWatchOnly:          i%2 == 0,
			AddrSchema:           nil,
		},
		TotalBalance: btcutil.Amount(1000 + i),
	}
}

func buildAccountResultValues(n int) []wallet.AccountResult {
	s := make([]wallet.AccountResult, n)
	for i := 0; i < n; i++ {
		s[i] = makeAccountResultValue(i)
	}
	return s
}

func buildAccountResultPointers(n int) []*wallet.AccountResult {
	s := make([]*wallet.AccountResult, n)
	for i := 0; i < n; i++ {
		s[i] = makeAccountResultPointer(i)
	}
	return s
}

// BenchmarkAccountResult_SliceBuild benchmarks building slices of AccountResult values vs pointers
func BenchmarkAccountResult_SliceBuild(b *testing.B) {
	// Reuse outpoint-style single-dimension datasets (count growth only).
	datasets := generateOutPointDatasets(outpointBenchConfig{
		outpointGrowth: scaleGrowth(8, exponentialGrowth()),
		iterations:     8,
	})
	for _, d := range datasets {
		prefix := fmt.Sprintf("%s-Accounts", d.name())
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var s []wallet.AccountResult
			for b.Loop() {
				s = buildAccountResultValues(d.numOutPoints)
			}
			sinkInt = len(s)
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var s []*wallet.AccountResult
			for b.Loop() {
				s = buildAccountResultPointers(d.numOutPoints)
			}
			sinkInt = len(s)
		})
	}
}

// BenchmarkAccountResult_SliceIterate benchmarks iterating over slices of AccountResult values vs pointers
func BenchmarkAccountResult_SliceIterate(b *testing.B) {
	datasets := generateOutPointDatasets(outpointBenchConfig{
		outpointGrowth: scaleGrowth(8, exponentialGrowth()),
		iterations:     8,
	})
	for _, d := range datasets {
		vals := buildAccountResultValues(d.numOutPoints)
		ptrs := buildAccountResultPointers(d.numOutPoints)
		prefix := fmt.Sprintf("%s-Accounts", d.name())
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				for _, val := range vals {
					acc += int64(val.TotalBalance)
					acc += int64(val.AccountNumber)
					acc += int64(val.ExternalKeyCount)
					acc += int64(val.InternalKeyCount)
					acc += int64(val.ImportedKeyCount)
					acc += int64(val.MasterKeyFingerprint)
					acc += int64(len(val.AccountName))
					if val.IsWatchOnly {
						acc++
					}
				}
			}
			sinkI64 = acc
		})
		b.Run(prefix+"/1-Pointers", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				for _, ptr := range ptrs {
					acc += int64(ptr.TotalBalance)
					acc += int64(ptr.AccountNumber)
					acc += int64(ptr.ExternalKeyCount)
					acc += int64(ptr.InternalKeyCount)
					acc += int64(ptr.ImportedKeyCount)
					acc += int64(ptr.MasterKeyFingerprint)
					acc += int64(len(ptr.AccountName))
					if ptr.IsWatchOnly {
						acc++
					}
				}
			}
			sinkI64 = acc
		})
	}
}

// BenchmarkAccountResult_SliceBuildAndIterate benchmarks building and iterating over slices
// of AccountResult values vs pointers with repeated reads.
func BenchmarkAccountResult_SliceBuildAndIterate(b *testing.B) {
	// Fixed small dataset for build+iterate stress with repeated reads.
	type ds struct{ n, max int }
	d := ds{n: 256, max: 256}
	for i := range 10 {
		prefix := fmt.Sprintf("%0*d-AccountsNReads%d", len(fmt.Sprintf("%d", d.max)), d.n, i)
		b.Run(prefix+"/0-Values", func(b *testing.B) {
			b.ReportAllocs()
			var acc int64
			for b.Loop() {
				vals := buildAccountResultValues(d.n)
				for range i {
					for _, val := range vals {
						acc += int64(val.TotalBalance)
						acc += int64(val.AccountNumber)
						acc += int64(val.ExternalKeyCount)
						acc += int64(val.InternalKeyCount)
						acc += int64(val.ImportedKeyCount)
						acc += int64(val.MasterKeyFingerprint)
						acc += int64(len(val.AccountName))
						if val.IsWatchOnly {
							acc++
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
				ptrs := buildAccountResultPointers(d.n)
				for range i {
					for _, ptr := range ptrs {
						acc += int64(ptr.TotalBalance)
						acc += int64(ptr.AccountNumber)
						acc += int64(ptr.ExternalKeyCount)
						acc += int64(ptr.InternalKeyCount)
						acc += int64(ptr.ImportedKeyCount)
						acc += int64(ptr.MasterKeyFingerprint)
						acc += int64(len(ptr.AccountName))
						if ptr.IsWatchOnly {
							acc++
						}
					}
				}
			}
			sinkI64 = acc
		})
	}
}
