package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	RedisAddr string `yaml:"redis_addr"`
	RedisDB   int    `yaml:"redis_db"`
	FloodN    int    `yaml:"flood_n"`
	FloodK    int    `yaml:"flood_k"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	err = decode(f, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg.RedisAddr == "" {
		return fmt.Errorf("redis_addr must not be empty")
	}
	if cfg.FloodN <= 0 {
		return fmt.Errorf("flood_n must be greater than 0")
	}
	if cfg.FloodK <= 0 {
		return fmt.Errorf("flood_k must be greater than 0")
	}
	return nil
}

func decode(r io.Reader, v interface{}) error {
	decoder := yaml.NewDecoder(r)
	return decoder.Decode(v)
}
