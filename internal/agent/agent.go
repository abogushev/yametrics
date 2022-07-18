package agent

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"go.uber.org/zap"
)

type Metrics struct {
	*runtime.MemStats
	PollCount   int64
	RandomValue float64
}

var logger *zap.SugaredLogger

func Run(ctx context.Context, l *zap.SugaredLogger) {
	mtrcs := &Metrics{&runtime.MemStats{}, 0, 0.0}
	logger = l
	go updateMetricsWithInterval(mtrcs, ctx)
	go sendMetricsWithInterval(mtrcs, ctx)
	logger.Info("agent started")
}

func schedule(f func(), ctx context.Context, duration time.Duration, name string) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			f()

		case <-ctx.Done():
			ticker.Stop()
			logger.Infof("cancel task: %s", name)
			return
		}
	}
}

func updateMetricsWithInterval(m *Metrics, ctx context.Context) {
	schedule(
		func() {
			runtime.ReadMemStats(m.MemStats)
			m.PollCount++
			m.RandomValue = rand.Float64()
		},
		ctx,
		2*time.Second,
		"collecting metrics")
}

func sendMetricsWithInterval(m *Metrics, ctx context.Context) {
	schedule(func() { sendMetrics(m) }, ctx, 10*time.Second, "sending metrics")
}

func sendMetrics(m *Metrics) {
	client := http.Client{}
	urls := mkURL(m)

	for i := 0; i < len(urls); i++ {
		r, err := client.Post(urls[i], "text/plain", nil)
		if err != nil {
			logger.Errorf("error in send metric: %s", err)
			fmt.Println(err)
		} else {
			r.Body.Close()
		}
	}
}

func mkURL(m *Metrics) []string {
	m2v := make(map[string]interface{})
	m2v["Alloc"] = m.Alloc
	m2v["BuckHashSys"] = m.BuckHashSys
	m2v["Frees"] = m.Frees
	m2v["GCCPUFraction"] = m.GCCPUFraction
	m2v["GCSys"] = m.GCSys
	m2v["HeapAlloc"] = m.HeapAlloc
	m2v["HeapIdle"] = m.HeapIdle
	m2v["HeapInuse"] = m.HeapInuse
	m2v["HeapObjects"] = m.HeapObjects
	m2v["HeapReleased"] = m.HeapReleased
	m2v["HeapSys"] = m.HeapSys
	m2v["LastGC"] = m.LastGC
	m2v["Lookups"] = m.Lookups
	m2v["MCacheInuse"] = m.MCacheInuse
	m2v["MCacheSys"] = m.MCacheSys
	m2v["MSpanInuse"] = m.MSpanInuse
	m2v["MSpanSys"] = m.MSpanSys
	m2v["Mallocs"] = m.Mallocs
	m2v["NextGC"] = m.NextGC
	m2v["NumForcedGC"] = m.NumForcedGC
	m2v["NumGC"] = m.NumGC
	m2v["OtherSys"] = m.OtherSys
	m2v["PauseTotalNs"] = m.PauseTotalNs
	m2v["StackInuse"] = m.StackInuse
	m2v["StackSys"] = m.StackSys
	m2v["Sys"] = m.Sys
	m2v["TotalAlloc"] = m.TotalAlloc
	m2v["RandomValue"] = m.RandomValue

	arr := make([]string, len(m2v)+1)
	i := 0
	for key, v := range m2v {
		arr[i] = fmt.Sprintf("http://127.0.0.1:8080/update/gauge/%s/%v", key, v)
		i++
	}
	arr[len(arr)-1] = fmt.Sprintf("http://127.0.0.1:8080/update/counter/%s/%v", "PollCount", m.PollCount)

	return arr
}