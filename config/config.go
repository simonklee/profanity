// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"github.com/BurntSushi/toml"
	"github.com/simonz05/profanity/types"
)

type Config struct {
	Listen string
	Region string
	Filter types.FilterType
	Redis  RedisConfig
}

type RedisConfig struct {
	DSN string `toml:"dsn"`
}

func ReadFile(filename string) (*Config, error) {
	config := new(Config)
	_, err := toml.DecodeFile(filename, config)
	return config, err
}

func ReadFileOrDefault(filename string) (*Config, error) {
	config := new(Config)
	_, err := toml.DecodeFile(filename, config)
	if err == nil {
		return config, nil
	}
	return &Config{}, nil
}
