package agent

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mkURL(t *testing.T) {
	m := &Metrics{&runtime.MemStats{}, 1, 2.0}
	url := "http://127.0.0.1:8080/update/gauge"
	runtime.ReadMemStats(m.MemStats)
	tests := []struct {
		name string
		m    *Metrics
		want []string
	}{
		{
			"correct urls",
			m,
			[]string{
				fmt.Sprintf("%v/Alloc/%v", url, m.Alloc),
				fmt.Sprintf("%v/BuckHashSys/%v", url, m.BuckHashSys),
				fmt.Sprintf("%v/Frees/%v", url, m.Frees),
				fmt.Sprintf("%v/GCCPUFraction/%v", url, m.GCCPUFraction),
				fmt.Sprintf("%v/GCSys/%v", url, m.GCSys),
				fmt.Sprintf("%v/HeapAlloc/%v", url, m.HeapAlloc),
				fmt.Sprintf("%v/HeapIdle/%v", url, m.HeapIdle),
				fmt.Sprintf("%v/HeapInuse/%v", url, m.HeapInuse),
				fmt.Sprintf("%v/HeapObjects/%v", url, m.HeapObjects),
				fmt.Sprintf("%v/HeapReleased/%v", url, m.HeapReleased),
				fmt.Sprintf("%v/HeapSys/%v", url, m.HeapSys),
				fmt.Sprintf("%v/LastGC/%v", url, m.LastGC),
				fmt.Sprintf("%v/Lookups/%v", url, m.Lookups),
				fmt.Sprintf("%v/MCacheInuse/%v", url, m.MCacheInuse),
				fmt.Sprintf("%v/MCacheSys/%v", url, m.MCacheSys),
				fmt.Sprintf("%v/MSpanInuse/%v", url, m.MSpanInuse),
				fmt.Sprintf("%v/MSpanSys/%v", url, m.MSpanSys),
				fmt.Sprintf("%v/Mallocs/%v", url, m.Mallocs),
				fmt.Sprintf("%v/NextGC/%v", url, m.NextGC),
				fmt.Sprintf("%v/NumForcedGC/%v", url, m.NumForcedGC),
				fmt.Sprintf("%v/NumGC/%v", url, m.NumGC),
				fmt.Sprintf("%v/OtherSys/%v", url, m.OtherSys),
				fmt.Sprintf("%v/PauseTotalNs/%v", url, m.PauseTotalNs),
				fmt.Sprintf("%v/StackInuse/%v", url, m.StackInuse),
				fmt.Sprintf("%v/StackSys/%v", url, m.StackSys),
				fmt.Sprintf("%v/Sys/%v", url, m.Sys),
				fmt.Sprintf("%v/TotalAlloc/%v", url, m.TotalAlloc),
				fmt.Sprintf("%v/RandomValue/%v", url, m.RandomValue),
				fmt.Sprintf("http://127.0.0.1:8080/update/counter/PollCount/%v", m.PollCount),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mkURL(tt.m); !assert.ElementsMatch(t, got, tt.want) {
				t.Errorf("mkURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
