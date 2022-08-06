package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Address string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	SignKey string `env:"KEY"`
}

func (cfg *ServerConfig) DefineFlags() {
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "server host:port")
	flag.StringVar(&cfg.SignKey, "k", "", "sign key")
}

func (cfg *ServerConfig) LoadFromEnv() {
	setIfDefined := func(key string, setV func(v string)) {
		if v, ok := os.LookupEnv(key); ok {
			setV(v)
		}
	}

	setIfDefined("ADDRESS", func(v string) { cfg.Address = v })
	setIfDefined("KEY", func(v string) { cfg.SignKey = v })
}
