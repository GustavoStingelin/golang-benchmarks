package pv

import "fmt"

// sink variables prevent the compiler from optimizing away benchmarked work.
// Shared across all benchmarks in this package.
var (
	sinkI64 int64
	sinkInt int
)

// growthFunc controls parameter scaling across iterations.
type growthFunc func(i int) int

func linearGrowth(step int) growthFunc {
	if step <= 0 {
		step = 1
	}
	return func(i int) int { return step * (i + 1) }
}

func exponentialGrowth() growthFunc {
	return func(i int) int { return 1 << i }
}

func constantGrowth(v int) growthFunc {
	return func(i int) int { return v }
}

// scaleGrowth multiplies a growth function by a constant base.
func scaleGrowth(base int, g growthFunc) growthFunc {
	if base <= 0 {
		base = 1
	}
	return func(i int) int { return base * g(i) }
}

type utxoDataset struct {
	numUtxos   int
	scriptSize int
	maxUtxos   int
	maxScript  int
}

// name returns a prefix used for b.Run to group runs in tools like vizb/benchstat.
// Example: "04096-Utxos-0034-Script".
func (d utxoDataset) name() string {
	uDigits := len(fmt.Sprintf("%d", d.maxUtxos))
	sDigits := len(fmt.Sprintf("%d", d.maxScript))
	return fmt.Sprintf("%0*d-Utxos-%0*d-Script", uDigits, d.numUtxos, sDigits, d.scriptSize)
}

type utxoBenchConfig struct {
	utxoGrowth   growthFunc
	scriptGrowth growthFunc
	iterations   int
}

func generateUtxoDatasets(c utxoBenchConfig) []utxoDataset {
	if c.iterations < 1 {
		c.iterations = 1
	}
	maxU := c.utxoGrowth(c.iterations - 1)
	maxS := c.scriptGrowth(c.iterations - 1)
	out := make([]utxoDataset, 0, c.iterations)
	for i := 0; i < c.iterations; i++ {
		out = append(out, utxoDataset{
			numUtxos:   c.utxoGrowth(i),
			scriptSize: c.scriptGrowth(i),
			maxUtxos:   maxU,
			maxScript:  maxS,
		})
	}
	return out
}

type outpointDataset struct {
	numOutPoints int
	maxOutPoints int
}

// name returns a prefix used for b.Run to group runs for OutPoint.
// Example: "04096-OutPoints".
func (d outpointDataset) name() string {
	uDigits := len(fmt.Sprintf("%d", d.maxOutPoints))
	return fmt.Sprintf("%0*d-OutPoints", uDigits, d.numOutPoints)
}

type outpointBenchConfig struct {
	outpointGrowth growthFunc
	iterations     int
}

func generateOutPointDatasets(c outpointBenchConfig) []outpointDataset {
	if c.iterations < 1 {
		c.iterations = 1
	}
	maxU := c.outpointGrowth(c.iterations - 1)
	out := make([]outpointDataset, 0, c.iterations)
	for i := 0; i < c.iterations; i++ {
		out = append(out, outpointDataset{
			numOutPoints: c.outpointGrowth(i),
			maxOutPoints: maxU,
		})
	}
	return out
}

// ---- Additional dataset helpers for TxIn, TxOut, and MsgTx ----

type txinDataset struct {
    numTxIns   int
    scriptSize int
    maxTxIns   int
    maxScript  int
}

// name returns a prefix for TxIn dataset. Example: "04096-TxIns-0034-Sig".
func (d txinDataset) name() string {
    uDigits := len(fmt.Sprintf("%d", d.maxTxIns))
    sDigits := len(fmt.Sprintf("%d", d.maxScript))
    return fmt.Sprintf("%0*d-TxIns-%0*d-Sig", uDigits, d.numTxIns, sDigits, d.scriptSize)
}

type txinBenchConfig struct {
    txinGrowth   growthFunc
    scriptGrowth growthFunc
    iterations   int
}

func generateTxInDatasets(c txinBenchConfig) []txinDataset {
    if c.iterations < 1 {
        c.iterations = 1
    }
    maxU := c.txinGrowth(c.iterations - 1)
    maxS := c.scriptGrowth(c.iterations - 1)
    out := make([]txinDataset, 0, c.iterations)
    for i := 0; i < c.iterations; i++ {
        out = append(out, txinDataset{
            numTxIns:   c.txinGrowth(i),
            scriptSize: c.scriptGrowth(i),
            maxTxIns:   maxU,
            maxScript:  maxS,
        })
    }
    return out
}

type txoutDataset struct {
    numTxOuts  int
    scriptSize int
    maxTxOuts  int
    maxScript  int
}

// name returns a prefix for TxOut dataset. Example: "04096-TxOuts-0034-Script".
func (d txoutDataset) name() string {
    uDigits := len(fmt.Sprintf("%d", d.maxTxOuts))
    sDigits := len(fmt.Sprintf("%d", d.maxScript))
    return fmt.Sprintf("%0*d-TxOuts-%0*d-Script", uDigits, d.numTxOuts, sDigits, d.scriptSize)
}

type txoutBenchConfig struct {
    txoutGrowth  growthFunc
    scriptGrowth growthFunc
    iterations   int
}

func generateTxOutDatasets(c txoutBenchConfig) []txoutDataset {
    if c.iterations < 1 {
        c.iterations = 1
    }
    maxU := c.txoutGrowth(c.iterations - 1)
    maxS := c.scriptGrowth(c.iterations - 1)
    out := make([]txoutDataset, 0, c.iterations)
    for i := 0; i < c.iterations; i++ {
        out = append(out, txoutDataset{
            numTxOuts:  c.txoutGrowth(i),
            scriptSize: c.scriptGrowth(i),
            maxTxOuts:  maxU,
            maxScript:  maxS,
        })
    }
    return out
}

type msgtxDataset struct {
    numTxs     int
    scriptSize int
    nInputs    int
    nOutputs   int
    maxTxs     int
    maxScript  int
}

// name returns a prefix for MsgTx dataset.
// Example: "04096-Txs-0034-Script-2x2" where 2x2 are in/out counts.
func (d msgtxDataset) name() string {
    uDigits := len(fmt.Sprintf("%d", d.maxTxs))
    sDigits := len(fmt.Sprintf("%d", d.maxScript))
    return fmt.Sprintf("%0*d-Txs-%0*d-Script-%dx%d", uDigits, d.numTxs, sDigits, d.scriptSize, d.nInputs, d.nOutputs)
}

type msgtxBenchConfig struct {
    txGrowth     growthFunc
    scriptGrowth growthFunc
    iterations   int
    nInputs      int
    nOutputs     int
}

func generateMsgTxDatasets(c msgtxBenchConfig) []msgtxDataset {
    if c.iterations < 1 {
        c.iterations = 1
    }
    if c.nInputs <= 0 {
        c.nInputs = 2
    }
    if c.nOutputs <= 0 {
        c.nOutputs = 2
    }
    maxU := c.txGrowth(c.iterations - 1)
    maxS := c.scriptGrowth(c.iterations - 1)
    out := make([]msgtxDataset, 0, c.iterations)
    for i := 0; i < c.iterations; i++ {
        out = append(out, msgtxDataset{
            numTxs:     c.txGrowth(i),
            scriptSize: c.scriptGrowth(i),
            nInputs:    c.nInputs,
            nOutputs:   c.nOutputs,
            maxTxs:     maxU,
            maxScript:  maxS,
        })
    }
    return out
}
