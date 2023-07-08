package config

import "os"

type Config interface {
	SECRET_KEY() string
}

type EnvConfig struct {
	secretKey string
}

func getEnvOrPanic(key string) string {
	value, set := os.LookupEnv(key)
	if !set {
		panic("Environment variable " + key + " is not set")
	}
	return value
}

func NewConfig() Config {
	return &EnvConfig{
		secretKey: getEnvOrPanic("OWL_SECRET_KEY"),
	}
}

func (c *EnvConfig) SECRET_KEY() string {
	return c.secretKey
}
