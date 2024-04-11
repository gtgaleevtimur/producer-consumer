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
	Address   string `env:"KAFKA_PRODUCER_ADDRESSES"`
	Addresses []string
}

func NewConfig() *Config {
	once.Do(func() {
		config = &Config{
			Address:   "localhost:9092",
			Addresses: []string{"localhost:9092"},
		}
		env.Parse(config)
		config.Addresses = strings.Split(config.Address, ",")
	})
	return config
}
