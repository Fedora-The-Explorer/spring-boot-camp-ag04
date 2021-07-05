package config

import "github.com/kelseyhightower/envconfig"

var Cfg Config

type Config struct {
	Api    apiConfig    `split_words:"true"`
}

type apiConfig struct {
	ReadWriteTimeoutMs int `split_words:"true" default:"10000"`
	Port               int `split_words:"true" default:"8080"`
}

func Load() {
	err := envconfig.Process("", &Cfg)
	if err != nil {
		panic(err)
	}
}