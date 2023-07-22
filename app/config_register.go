package app

type ConfigRegister struct {
	configs map[string]interface{}
}

type RegisteredConfig struct {
	Name   string
	Config interface{}
}

func NewConfigRegister() *ConfigRegister {
	return &ConfigRegister{configs: map[string]interface{}{}}
}

func (r *ConfigRegister) Register(name string, config interface{}) {
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

func (r *ConfigRegister) GetConfig(name string) interface{} {
	return r.configs[name]
}
