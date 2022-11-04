package storage

import (
	"runtime"
	"testing"
)

var M = &Metrics{MemStats: &runtime.MemStats{}, PollCount: 0, RandomValue: 0.0}

func BenchmarkMetricToMap(b *testing.B) {
	b.Run("current", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			M.MetricToMaps()
		}
	})
	b.Run("old", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			M.MetricToMapsOld()
		}
	})

}
