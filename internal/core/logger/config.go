package logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type LogConfig struct {
	Level  string `envconfig:"LEVEL" default:"DEBUG"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

func NewConfig() (LogConfig, error) {
	var conf LogConfig

	if err := envconfig.Process("LOG", &conf); err != nil {
		return LogConfig{}, fmt.Errorf("error loading env: %v", err)
	}
	return conf, nil
}

func NewConfigMust() LogConfig {
	conf, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return conf
}
