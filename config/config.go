package config

import "os"

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
