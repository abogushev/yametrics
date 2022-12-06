package config

import (
	"flag"
	"os"
	"time"
	"yametrics/internal/configfile"
	"yametrics/internal/durationextension"
)

type AgentConfig struct {
	Address        string                     `env:"ADDRESS" envDefault:"127.0.0.1:8080" json:"address"`
	ReportInterval durationextension.Duration `env:"REPORT_INTERVAL" envDefault:"10s" json:"report_interval"`
	PollInterval   durationextension.Duration `env:"POLL_INTERVAL" envDefault:"2s" json:"poll_interval"`
	SignKey        string                     `env:"KEY"`
	CryptoKeyPath  string                     `env:"CRYPTO_KEY" json:"crypto_key"`
	configPath     string
}

func NewAgentConfig() *AgentConfig {
	cfg := &AgentConfig{}
	cfg.defineFlags()
	flag.Parse()
	cfg.readConfigFile()
	flag.Parse()
	cfg.loadFromEnv()
	return cfg
}

func (cfg *AgentConfig) defineFlags() {
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "server address")
	flag.DurationVar(&cfg.ReportInterval.Duration, "r", time.Second*10, "report interval")
	flag.DurationVar(&cfg.PollInterval.Duration, "p", time.Second*2, "poll interval")
	flag.StringVar(&cfg.SignKey, "k", "", "sign key")
	flag.StringVar(&cfg.CryptoKeyPath, "crypto-key", "public_key.pem", "path to public key")
	flag.StringVar(&cfg.configPath, "c", "", "path to config file")
}

func (cfg *AgentConfig) loadFromEnv() {
	setIfDefined := func(key string, setV func(v string)) {
		if v, ok := os.LookupEnv(key); ok {
			setV(v)
		}
	}
	setIfDefined("ADDRESS", func(v string) { cfg.Address = v })
	setIfDefined("REPORT_INTERVAL", func(v string) { cfg.ReportInterval.Duration, _ = time.ParseDuration(v) })
	setIfDefined("POLL_INTERVAL", func(v string) { cfg.PollInterval.Duration, _ = time.ParseDuration(v) })
	setIfDefined("KEY", func(v string) { cfg.SignKey = v })
	setIfDefined("CRYPTO_KEY", func(v string) { cfg.CryptoKeyPath = v })
}

func (cfg *AgentConfig) readConfigFile() {
	if cfg.configPath != "" {
		err := configfile.ReadConfig(cfg.configPath, cfg)
		if err != nil {
			panic(err)
		}
	}
}
