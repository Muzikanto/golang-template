package config

import (
	"fmt"
	"go-backend-template/src/utils/config"
	"go.uber.org/dig"
)

type Config struct {
	HttpHost          string `envconfig:"HTTP_HOST"`
	HttpPort          int    `envconfig:"HTTP_PORT"`
	HttpDetailedError bool   `envconfig:"HTTP_DETAILED_ERROR"`

	DatabaseURL string `envconfig:"DATABASE_URL"`

	AccessTokenExpiresTTL int    `envconfig:"ACCESS_TOKEN_EXPIRES_TTL"`
	AccessTokenSecret     string `envconfig:"ACCESS_TOKEN_SECRET"`
}

func CreateConfig(env string) *Config {
	var cfg = config.ParseEnvConfig(&Config{}, env)

	return cfg
}

func (r *Config) Provide(container *dig.Scope) *Config {
	config.ProvideConfig[Config](container, r)

	return r
}

// Http

func (r *Config) HTTP() *HttpConfig {
	return &HttpConfig{
		host:          r.HttpHost,
		port:          r.HttpPort,
		detailedError: r.HttpDetailedError,
	}
}

type HttpConfig struct {
	host          string
	port          int
	detailedError bool
}

func (c *HttpConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

// DB

func (r *Config) Database() *DatabaseConfig {
	return &DatabaseConfig{
		url: r.DatabaseURL,
	}
}

type DatabaseConfig struct {
	url string
}

func (c *DatabaseConfig) ConnString() string {
	return c.url
}
