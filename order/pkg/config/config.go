package config

import (
	"fmt"

	confutil "github.com/swarit-pandey/e-commerce/common/conf"
)

type RootConfig struct {
	ConfigPath string
}

func Default() *RootConfig {
	return &RootConfig{}
}

func (cfg *RootConfig) Load() (*OrderConfig, error) {
	if cfg.ConfigPath == "" {
		cfg.ConfigPath = "../conf/app.yaml"
	}

	configLoad := Config{}

	err := confutil.LoadAndUnmarhsal(cfg.ConfigPath, configLoad)
	if err != nil {
		return nil, fmt.Errorf("notification: %v", err)
	}

	return &configLoad.Order, nil
}
