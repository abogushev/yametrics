package storage

import (
	"runtime"
	"yametrics/internal/agent/models/api"
)

type Metrics struct {
	*runtime.MemStats
	PollCount   int64
	RandomValue float64
}

func (m *Metrics) ToAPI() []api.Metrics {
	result := make([]api.Metrics, 0)
	m.OperateOverMetricMaps(
		func(s string, f float64) {
			result = append(result, api.Metrics{ID: s, MType: api.GAUGE, Value: &f})
		},
		func(s string, i int64) {
			result = append(result, api.Metrics{ID: s, MType: api.COUNTER, Delta: &i})
		},
	)
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

	m2counter := make(map[string]int64)
	m2counter["PollCount"] = m.PollCount

	return m2gauge, m2counter
}