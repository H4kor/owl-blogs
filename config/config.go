package config

import "os"

const (
	OWL_VERSION       = "0.3.6"
	SITE_CONFIG       = "site_config"
	ACT_PUB_CONF_NAME = "activity_pub"
)

type Config interface {
}

type EnvConfig struct {
}

func getEnvOrPanic(key string) string {
	value, set := os.LookupEnv(key)
	if !set {
		panic("Environment variable " + key + " is not set")
	}
	return value
}

func NewConfig() Config {
	return &EnvConfig{}
}
