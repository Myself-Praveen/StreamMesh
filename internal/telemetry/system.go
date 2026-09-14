package telemetry

import (
	"runtime"
)

// SystemMetrics holds OS and Go runtime metrics
type SystemMetrics struct {
	NumGoroutines int
	AllocBytes    uint64
	SysBytes      uint64
	NumGC         uint32
}

// GetSystemMetrics returns a snapshot of current system metrics
func GetSystemMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return SystemMetrics{
		NumGoroutines: runtime.NumGoroutine(),
		AllocBytes:    m.Alloc,
		SysBytes:      m.Sys,
		NumGC:         m.NumGC,
	}
}
