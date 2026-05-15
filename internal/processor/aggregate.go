package processor

import (
	"fmt"
	"sort"
	"sync"

	"github.com/yourorg/logslice/internal/parser"
)

// AggregateMode defines how values are aggregated per group key.
type AggregateMode string

const (
	AggCount AggregateMode = "count"
	AggSum   AggregateMode = "sum"
	AggMin   AggregateMode = "min"
	AggMax   AggregateMode = "max"
)

// Aggregator groups log lines by a field value and computes an aggregate
// over a numeric target field (or line count for "count" mode).
type Aggregator struct {
	groupField  string
	valueField  string
	mode        AggregateMode
	mu          sync.Mutex
	buckets     map[string]float64
	counts      map[string]int64
}

// NewAggregator creates an Aggregator. groupField is the field to group by,
// valueField is the numeric field to aggregate (ignored for count mode),
// and mode is one of: count, sum, min, max.
func NewAggregator(groupField, valueField string, mode AggregateMode) (*Aggregator, error) {
	if groupField == "" {
		return nil, fmt.Errorf("aggregator: groupField must not be empty")
	}
	switch mode {
	case AggCount, AggSum, AggMin, AggMax:
		// valid
	default:
		return nil, fmt.Errorf("aggregator: unknown mode %q", mode)
	}
	if mode != AggCount && valueField == "" {
		return nil, fmt.Errorf("aggregator: valueField required for mode %q", mode)
	}
	return &Aggregator{
		groupField: groupField,
		valueField: valueField,
		mode:       mode,
		buckets:    make(map[string]float64),
		counts:     make(map[string]int64),
	}, nil
}

// Ingest processes a single parsed log line into the aggregation buckets.
func (a *Aggregator) Ingest(line parser.LogLine) {
	key, ok := line.Fields[a.groupField]
	if !ok {
		key = "(none)"
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.counts[key]++
	if a.mode == AggCount {
		return
	}
	var val float64
	if raw, exists := line.Fields[a.valueField]; exists {
		fmt.Sscanf(raw, "%f", &val)
	}
	switch a.mode {
	case AggSum:
		a.buckets[key] += val
	case AggMin:
		if _, seen := a.buckets[key]; !seen || val < a.buckets[key] {
			a.buckets[key] = val
		}
	case AggMax:
		if _, seen := a.buckets[key]; !seen || val > a.buckets[key] {
			a.buckets[key] = val
		}
	}
}

// AggregateResult holds one aggregated bucket result.
type AggregateResult struct {
	Key   string
	Value float64
	Count int64
}

// Results returns sorted aggregate results (ascending by key).
func (a *Aggregator) Results() []AggregateResult {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]AggregateResult, 0, len(a.counts))
	for k, c := range a.counts {
		v := a.buckets[k]
		if a.mode == AggCount {
			v = float64(c)
		}
		out = append(out, AggregateResult{Key: k, Value: v, Count: c})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// Reset clears all accumulated state.
func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.buckets = make(map[string]float64)
	a.counts = make(map[string]int64)
}
