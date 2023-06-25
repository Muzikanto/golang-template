package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
	"go.uber.org/dig"
	"log"
)

func ParseEnvConfig[C interface{}](config *C, envPath string) *C {
	if envPath != "" {
		if err := gotenv.OverLoad(envPath); err != nil {
			log.Fatal(err)
		}
	}

	if err := envconfig.Process("", config); err != nil {
		log.Println("Invalid config")
		log.Fatal(err)
	}

	return config
}

func ProvideConfig[C interface{}](container *dig.Scope, config *C) {
	if err := container.Provide(func() *C {
		return config
	}); err != nil {
		log.Fatal(err)
	}
}
