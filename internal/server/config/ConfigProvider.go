package config

import "flag"

type ConfigProvider struct {
	StorageCfg *MetricsStorageConfig
	ServerCfg  *ServerConfig
}

func NewConfigProvider() *ConfigProvider {
	storageCfg := &MetricsStorageConfig{}
	serverCfg := &ServerConfig{}

	storageCfg.DefineFlags()
	serverCfg.DefineFlags()

	flag.Parse()

	storageCfg.LoadFromEnv()
	serverCfg.LoadFromEnv()

	return &ConfigProvider{StorageCfg: storageCfg, ServerCfg: serverCfg}
}
