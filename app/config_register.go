package app

import "owl-blogs/domain/model"

type AppConfig interface {
	model.Formable
}

type ConfigRegister struct {
	configs map[string]AppConfig
}

type RegisteredConfig struct {
	Name   string
	Config AppConfig
}

func NewConfigRegister() *ConfigRegister {
	return &ConfigRegister{configs: map[string]AppConfig{}}
}

func (r *ConfigRegister) Register(name string, config AppConfig) {
	r.configs[name] = config
}

func (r *ConfigRegister) Configs() []RegisteredConfig {
	var configs []RegisteredConfig
	for name, config := range r.configs {
		configs = append(configs, RegisteredConfig{
			Name:   name,
			Config: config,
		})
	}
	return configs
}

func (r *ConfigRegister) GetConfig(name string) AppConfig {
	return r.configs[name]
}
