package storage

import (
	"runtime"
	"yametrics/internal/metricscrypto"
	"yametrics/internal/protocol"
)

type Metrics struct {
	*runtime.MemStats
	PollCount       int64
	RandomValue     float64
	TotalMemory     float64
	FreeMemory      float64
	CPUutilization1 float64
	m2gauge         map[string]float64
	m2counter       map[string]int64
}

func NewMetrics() *Metrics {
	return &Metrics{MemStats: &runtime.MemStats{}, PollCount: 0, RandomValue: 0.0, m2gauge: map[string]float64{}, m2counter: map[string]int64{}}
}

func (m *Metrics) ToAPI() []protocol.Metrics {
	result := make([]protocol.Metrics, 0)
	m.OperateOverMetricMaps(
		func(s string, f float64) {
			result = append(result, protocol.Metrics{ID: s, MType: protocol.GAUGE, Value: &f})
		},
		func(s string, i int64) {
			result = append(result, protocol.Metrics{ID: s, MType: protocol.COUNTER, Delta: &i})
		},
	)
	return result
}

func (m *Metrics) ToAPIWithSign(key string) []protocol.Metrics {
	result := m.ToAPI()
	for i := 0; i < len(result); i++ {
		result[i].Hash = metricscrypto.GetMetricSign(result[i], key)
	}
	return result
}

func (m *Metrics) OperateOverMetricMaps(gaugeF func(string, float64), countersF func(string, int64)) {
	gauges, counters := m.MetricToMaps()
	for k, v := range gauges {
		gaugeF(k, v)
	}
	for k, v := range counters {
		countersF(k, v)
	}
}

func (m *Metrics) MetricToMaps() (map[string]float64, map[string]int64) {
	m.m2gauge["Alloc"] = float64(m.Alloc)
	m.m2gauge["BuckHashSys"] = float64(m.BuckHashSys)
	m.m2gauge["Frees"] = float64(m.Frees)
	m.m2gauge["GCCPUFraction"] = m.GCCPUFraction
	m.m2gauge["GCSys"] = float64(m.GCSys)
	m.m2gauge["HeapAlloc"] = float64(m.HeapAlloc)
	m.m2gauge["HeapIdle"] = float64(m.HeapIdle)
	m.m2gauge["HeapInuse"] = float64(m.HeapInuse)
	m.m2gauge["HeapObjects"] = float64(m.HeapObjects)
	m.m2gauge["HeapReleased"] = float64(m.HeapReleased)
	m.m2gauge["HeapSys"] = float64(m.HeapSys)
	m.m2gauge["LastGC"] = float64(m.LastGC)
	m.m2gauge["Lookups"] = float64(m.Lookups)
	m.m2gauge["MCacheInuse"] = float64(m.MCacheInuse)
	m.m2gauge["MCacheSys"] = float64(m.MCacheSys)
	m.m2gauge["MSpanInuse"] = float64(m.MSpanInuse)
	m.m2gauge["MSpanSys"] = float64(m.MSpanSys)
	m.m2gauge["Mallocs"] = float64(m.Mallocs)
	m.m2gauge["NextGC"] = float64(m.NextGC)
	m.m2gauge["NumForcedGC"] = float64(m.NumForcedGC)
	m.m2gauge["NumGC"] = float64(m.NumGC)
	m.m2gauge["OtherSys"] = float64(m.OtherSys)
	m.m2gauge["PauseTotalNs"] = float64(m.PauseTotalNs)
	m.m2gauge["StackInuse"] = float64(m.StackInuse)
	m.m2gauge["StackSys"] = float64(m.StackSys)
	m.m2gauge["Sys"] = float64(m.Sys)
	m.m2gauge["TotalAlloc"] = float64(m.TotalAlloc)
	m.m2gauge["RandomValue"] = m.RandomValue
	m.m2gauge["TotalMemory"] = m.TotalMemory
	m.m2gauge["FreeMemory"] = m.FreeMemory
	m.m2gauge["CPUutilization1"] = m.CPUutilization1
	m.m2counter["PollCount"] = m.PollCount

	return m.m2gauge, m.m2counter
}

func (m *Metrics) MetricToMapsOld() (map[string]float64, map[string]int64) {
	m2gauge := make(map[string]float64)
	m2gauge["Alloc"] = float64(m.Alloc)
	m2gauge["BuckHashSys"] = float64(m.BuckHashSys)
	m2gauge["Frees"] = float64(m.Frees)
	m2gauge["GCCPUFraction"] = m.GCCPUFraction
	m2gauge["GCSys"] = float64(m.GCSys)
	m2gauge["HeapAlloc"] = float64(m.HeapAlloc)
	m2gauge["HeapIdle"] = float64(m.HeapIdle)
	m2gauge["HeapInuse"] = float64(m.HeapInuse)
	m2gauge["HeapObjects"] = float64(m.HeapObjects)
	m2gauge["HeapReleased"] = float64(m.HeapReleased)
	m2gauge["HeapSys"] = float64(m.HeapSys)
	m2gauge["LastGC"] = float64(m.LastGC)
	m2gauge["Lookups"] = float64(m.Lookups)
	m2gauge["MCacheInuse"] = float64(m.MCacheInuse)
	m2gauge["MCacheSys"] = float64(m.MCacheSys)
	m2gauge["MSpanInuse"] = float64(m.MSpanInuse)
	m2gauge["MSpanSys"] = float64(m.MSpanSys)
	m2gauge["Mallocs"] = float64(m.Mallocs)
	m2gauge["NextGC"] = float64(m.NextGC)
	m2gauge["NumForcedGC"] = float64(m.NumForcedGC)
	m2gauge["NumGC"] = float64(m.NumGC)
	m2gauge["OtherSys"] = float64(m.OtherSys)
	m2gauge["PauseTotalNs"] = float64(m.PauseTotalNs)
	m2gauge["StackInuse"] = float64(m.StackInuse)
	m2gauge["StackSys"] = float64(m.StackSys)
	m2gauge["Sys"] = float64(m.Sys)
	m2gauge["TotalAlloc"] = float64(m.TotalAlloc)
	m2gauge["RandomValue"] = m.RandomValue
	m2gauge["TotalMemory"] = m.TotalMemory
	m2gauge["FreeMemory"] = m.FreeMemory
	m2gauge["CPUutilization1"] = m.CPUutilization1

	m2counter := make(map[string]int64)
	m2counter["PollCount"] = m.PollCount

	return m2gauge, m2counter
}
