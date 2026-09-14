package telemetry

import (
	"sort"
	"sync"
)

// Histogram calculates percentiles over a rolling window of samples
type Histogram struct {
	mu      sync.Mutex
	samples []float64
	head    int
	count   int
	size    int
}

// NewHistogram creates a new Histogram with a given window size
func NewHistogram(size int) *Histogram {
	return &Histogram{
		samples: make([]float64, size),
		size:    size,
	}
}

// Record records a latency value
func (h *Histogram) Record(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.samples[h.head] = value
	h.head = (h.head + 1) % h.size
	if h.count < h.size {
		h.count++
	}
}

// Percentile returns the p-th percentile (e.g. 99 for p99)
func (h *Histogram) Percentile(p float64) float64 {
	h.mu.Lock()
	if h.count == 0 {
		h.mu.Unlock()
		return 0.0
	}
	
	// Create a copy to sort
	data := make([]float64, h.count)
	for i := 0; i < h.count; i++ {
		data[i] = h.samples[i]
	}
	h.mu.Unlock()

	sort.Float64s(data)
	
	idx := int(float64(h.count) * (p / 100.0))
	if idx >= h.count {
		idx = h.count - 1
	}
	
	return data[idx]
}

// Snapshots returns common percentiles
func (h *Histogram) Snapshots() map[string]float64 {
	return map[string]float64{
		"p50": h.Percentile(50),
		"p90": h.Percentile(90),
		"p99": h.Percentile(99),
	}
}

// GlobalLatencyHistogram is a singleton for global message latency
var GlobalLatencyHistogram = NewHistogram(10000)
