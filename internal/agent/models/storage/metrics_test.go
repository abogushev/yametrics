package storage

import (
	"testing"
)

var M = NewMetrics()

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
