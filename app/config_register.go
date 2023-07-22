package app

type ConfigRegister struct {
	configs map[string]interface{}
}

func NewConfigRegister() *ConfigRegister {
	return &ConfigRegister{configs: map[string]interface{}{}}
}

func (r *ConfigRegister) Register(name string, config interface{}) {
	r.configs[name] = config
}
