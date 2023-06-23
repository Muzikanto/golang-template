package config

import (
	"fmt"
	"go-backend-template/src/utils/database"
	"go-backend-template/src/utils/http"
	"go.uber.org/dig"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
)

// Config

type Config struct {
	HttpHost          string `envconfig:"HTTP_HOST"`
	HttpPort          int    `envconfig:"HTTP_PORT"`
	HttpDetailedError bool   `envconfig:"HTTP_DETAILED_ERROR"`

	DatabaseURL string `envconfig:"DATABASE_URL"`

	AccessTokenExpiresTTL int    `envconfig:"ACCESS_TOKEN_EXPIRES_TTL"`
	AccessTokenSecret     string `envconfig:"ACCESS_TOKEN_SECRET"`
}

func ParseEnv(container *dig.Container, envPath string) *Config {
	if envPath != "" {
		if err := gotenv.OverLoad(envPath); err != nil {
			log.Fatal(err)
		}
	}

	var config Config

	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}

	// di
	if err := container.Provide(func() *Config {
		return &config
	}); err != nil {
		log.Fatal(err)
	}

	return &config
}

func (c *Config) HTTP() http.Config {
	return &httpConfig{
		host:          c.HttpHost,
		port:          c.HttpPort,
		detailedError: c.HttpDetailedError,
	}
}

func (c *Config) Database() database.Config {
	return &databaseConfig{
		url: c.DatabaseURL,
	}
}

// HTTP

type httpConfig struct {
	host          string
	port          int
	detailedError bool
}

func (c *httpConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func (c *httpConfig) DetailedError() bool {
	return c.detailedError
}

// Database

type databaseConfig struct {
	url string
}

func (c *databaseConfig) ConnString() string {
	return c.url
}

// Auth

type authConfig struct {
	accessTokenExpiresTTL int
	accessTokenSecret     string
}

func (c *authConfig) AccessTokenSecret() string {
	return c.accessTokenSecret
}

func (c *authConfig) AccessTokenExpiresDate() time.Time {
	duration := time.Duration(c.accessTokenExpiresTTL)
	return time.Now().UTC().Add(time.Minute * duration)
}
