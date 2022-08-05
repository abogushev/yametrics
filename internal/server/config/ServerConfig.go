package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Address string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

func (cfg *ServerConfig) DefineFlags() {
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "server host:port")
}

func (cfg *ServerConfig) LoadFromEnv() {
	if address, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.Address = address
	}
}
