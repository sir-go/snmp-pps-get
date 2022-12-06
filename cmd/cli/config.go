package main

import (
	"strings"

	"github.com/go-yaml/yaml"
)

const ConfigPath = "config.yaml"

type (
	cfgOIDs struct {
		Models []string `yaml:"models"`
		Tx     string   `yaml:"tx"`
		Rx     string   `yaml:"rx"`
	}

	Config struct {
		Oids    []cfgOIDs `yaml:"oids"`
		Default struct {
			Tx string `yaml:"tx"`
			Rx string `yaml:"rx"`
		} `yaml:"default"`
		SnmpCommunity string `yaml:"snmp-community"`
	}
)

func getOids(swModel string, cfg *Config) (oidTx string, oidRx string) {
	for _, oidCfg := range cfg.Oids {
		for _, model := range oidCfg.Models {
			if strings.Contains(swModel, model) {
				return oidCfg.Tx, oidCfg.Rx
			}
		}
	}
	return cfg.Default.Tx, cfg.Default.Rx
}

func LoadConfig(b []byte) (cfg *Config, err error) {
	cfg = new(Config)
	if err = yaml.UnmarshalStrict(b, cfg); err != nil {
		return nil, err
	}
	return
}
