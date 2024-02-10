package config

import (
	"github.com/jessevdk/go-flags"
	configutil "github.com/nihcet/go-lib/pkg/util/env"
)

type Config struct {
	configutil.ServerConfig
}

var config Config

func LoadConfig() {
	_, err := flags.Parse(&config)
	if err != nil {
		panic(err)
	}

}

func Get() Config {
	return config
}
