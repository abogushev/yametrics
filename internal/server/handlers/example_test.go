package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики
	Value *float64 `json:"value,omitempty"` // значение метрики
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func sendMultipleMetricsV2(metrcis []Metrics) {
	body, err := json.Marshal(metrcis)
	if err != nil {
		log.Fatalf("error on  Marshal metric: %v", err)
	}
	if r, err := http.Post("http://example.com/updates", "application/json", bytes.NewBuffer(body)); err != nil {
		log.Fatalf("error in send metric: %v", err)
	} else {
		r.Body.Close()
	}

	log.Printf("send successful")
}

func sendMetricsV2(metrcis []Metrics) {
	for i := 0; i < len(metrcis); i++ {
		json, err := json.Marshal(metrcis[i])
		if err != nil {
			log.Fatalf("error in Marshal metric: %v", err)
		}
		r, err := http.Post("http://example.com/update", "application/json", bytes.NewBuffer(json))
		if err != nil {
			log.Fatalf("error in send metric: %v", err)
		}
		r.Body.Close()
	}
	log.Printf("send successful")
}

func sendMetricsV1(gauges map[string]float64, counters map[string]int64) {
	rqsts := make([]string, 0, len(gauges)+len(counters))

	for key, v := range gauges {
		rqsts = append(rqsts, fmt.Sprintf("http://example.com/update/gauge/%v/%v", key, v))
	}
	for key, v := range counters {
		rqsts = append(rqsts, fmt.Sprintf("http://example.com/update/counter/%v/%v", key, v))
	}
	for _, rq := range rqsts {
		r, err := http.Post(rq, "text/plain", nil)
		if err != nil {
			log.Fatal(err)
		}
		r.Body.Close()
	}
	log.Printf("send successful")
}

func Example() {
	//example how to send single metric
	sendMetricsV1(make(map[string]float64), make(map[string]int64))
	//send single metric in json format
	sendMetricsV2(make([]Metrics, 0))
	//send multipe metric in json format
	sendMultipleMetricsV2(make([]Metrics, 0))
}
