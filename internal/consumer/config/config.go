package config

import (
	"strings"
	"sync"

	"github.com/caarlos0/env"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	Address           string `env:"KAFKA_PRODUCER_ADDRESSES"`
	Addresses         []string
	PrometheusAddress string `env:"PROMETHEUS_ADDRESS"`
	Interval          int    `env:"INTERVAL"`
	Count             int    `env:"COUNT"`
}

func NewConfig() *Config {
	once.Do(func() {
		config = &Config{
			Address:           "localhost:9093",
			Addresses:         []string{"localhost:9093"},
			PrometheusAddress: ":8080",
			Interval:          10,
			Count:             10,
		}
		env.Parse(config)
		config.Addresses = strings.Split(config.Address, ",")
	})
	return config
}
