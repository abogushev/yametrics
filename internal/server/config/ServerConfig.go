package config

import (
	"flag"
	"os"
	"strconv"
	"time"
	"yametrics/internal/configfile"
	"yametrics/internal/durationextension"
)

type ServerConfig struct {
	Address       string                     `env:"ADDRESS" envDefault:"127.0.0.1:8080" json:"address"`
	SignKey       string                     `env:"KEY"`
	CryptoKeyPath string                     `env:"CRYPTO_KEY" json:"crypto_key"`
	StoreInterval durationextension.Duration `env:"STORE_INTERVAL" envDefault:"300s" json:"store_interval"`
	StoreFile     string                     `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json" json:"store_file"`
	Restore       bool                       `env:"RESTORE" envDefault:"true" json:"restore"`
	DBURL         string                     `env:"DATABASE_DSN" json:"database_dsn"`
	TrustedSubnet string                     `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	configPath    string
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	cfg.defineFlags()
	flag.Parse()
	cfg.readConfigFile()
	flag.Parse()
	cfg.loadFromEnv()
	return cfg
}

func (cfg *ServerConfig) defineFlags() {
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "server host:port")
	flag.StringVar(&cfg.SignKey, "k", "", "sign key")
	flag.StringVar(&cfg.CryptoKeyPath, "crypto-key", "private_key.pem", "path to private key")
	flag.DurationVar(&cfg.StoreInterval.Duration, "i", time.Second*300, "save metrics interval")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "save metrics file")
	flag.BoolVar(&cfg.Restore, "r", true, "restore metrics ?")
	flag.StringVar(&cfg.DBURL, "d", "", "db connection url, exmlp: postgres://username:password@localhost:5432/database_name")
	flag.StringVar(&cfg.configPath, "c", "", "path to config file")
	flag.StringVar(&cfg.configPath, "t", "", "trusted subnet")
}

func (cfg *ServerConfig) loadFromEnv() {
	setIfDefined := func(key string, setV func(v string)) {
		if v, ok := os.LookupEnv(key); ok {
			setV(v)
		}
	}

	setIfDefined("ADDRESS", func(v string) { cfg.Address = v })
	setIfDefined("KEY", func(v string) { cfg.SignKey = v })
	setIfDefined("CRYPTO_KEY", func(v string) { cfg.CryptoKeyPath = v })
	setIfDefined("STORE_INTERVAL", func(v string) { cfg.StoreInterval.Duration, _ = time.ParseDuration(v) })
	setIfDefined("STORE_FILE", func(v string) { cfg.StoreFile = v })
	setIfDefined("RESTORE", func(v string) { cfg.Restore, _ = strconv.ParseBool(v) })
	setIfDefined("DATABASE_DSN", func(v string) { cfg.DBURL = v })
	setIfDefined("TRUSTED_SUBNET", func(v string) { cfg.TrustedSubnet = v })
}

func (cfg *ServerConfig) readConfigFile() {
	if cfg.configPath != "" {
		err := configfile.ReadConfig(cfg.configPath, cfg)
		if err != nil {
			panic(err)
		}
	}
}
