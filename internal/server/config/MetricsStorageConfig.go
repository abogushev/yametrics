package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type MetricsStorageConfig struct {
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func (cfg *MetricsStorageConfig) DefineFlags() {
	flag.DurationVar(&cfg.StoreInterval, "i", time.Second*300, "save metrics interval")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "save metrics file")
	flag.BoolVar(&cfg.Restore, "r", true, "restore metrics ?")
}

func (cfg *MetricsStorageConfig) LoadFromEnv() {
	setIfDefined := func(key string, setV func(v string)) {
		if v, ok := os.LookupEnv(key); ok {
			setV(v)
		}
	}
	setIfDefined("STORE_INTERVAL", func(v string) { cfg.StoreInterval, _ = time.ParseDuration(v) })
	setIfDefined("STORE_FILE", func(v string) { cfg.StoreFile = v })
	setIfDefined("RESTORE", func(v string) { cfg.Restore, _ = strconv.ParseBool(v) })
}