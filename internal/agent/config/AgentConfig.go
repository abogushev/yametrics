package config

import (
	"flag"
	"os"
	"time"
)

type AgentConfig struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	SignKey        string        `env:"KEY"`
}

func (cfg *AgentConfig) DefineFlags() {
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "server address")
	flag.DurationVar(&cfg.ReportInterval, "r", time.Second*10, "report interval")
	flag.DurationVar(&cfg.PollInterval, "p", time.Second*2, "poll interval")
	flag.StringVar(&cfg.SignKey, "k", "", "sign key")
}

func (cfg *AgentConfig) LoadFromEnv() {
	setIfDefined := func(key string, setV func(v string)) {
		if v, ok := os.LookupEnv(key); ok {
			setV(v)
		}
	}
	setIfDefined("ADDRESS", func(v string) { cfg.Address = v })
	setIfDefined("REPORT_INTERVAL", func(v string) { cfg.ReportInterval, _ = time.ParseDuration(v) })
	setIfDefined("POLL_INTERVAL", func(v string) { cfg.PollInterval, _ = time.ParseDuration(v) })
	setIfDefined("KEY", func(v string) { cfg.SignKey = v })
}
