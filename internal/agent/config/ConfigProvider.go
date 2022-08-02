package config

import "flag"

type ConfigProvider struct {
	AgentCfg *AgentConfig
}

func NewConfigProvider() *ConfigProvider {
	agentCfg := &AgentConfig{}

	agentCfg.DefineFlags()
	flag.Parse()

	agentCfg.LoadFromEnv()

	return &ConfigProvider{AgentCfg: agentCfg}
}
