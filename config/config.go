package config

import (
	"fmt"

	"github.com/Netflix/go-env"
)

type Config struct {
	JaegerGRPCEndpoint string `env:"JAEGER_GRPC_ENDPOINT"`
}

func Get() *Config {
	var cfg Config
	envar, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		panic(err)
	}

	if err := env.Unmarshal(envar, &cfg); err != nil {
		panic(err)
	}

	fmt.Println("JAEGER_GRPC_ENDPOINT", cfg.JaegerGRPCEndpoint)

	return &cfg
}
