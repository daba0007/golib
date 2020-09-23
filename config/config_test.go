package config

import (
	"reflect"
	"testing"
)

// DBConfig database config
type DBConfig struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Address  string `json:"address" yaml:"address"`
	DBName   string `json:"dbName" yaml:"dbName"`
}

// ESConfig elasticsearch config
type ESConfig struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Address  string `json:"address" yaml:"address"`
}

// GwayConfig gway config
type Gateway struct {
	Address string `json:"address" yaml:"address"`
	Count   int    `json:"count" yaml:"count"`
}

// RabbitmqConfig Rabbitmq config
type RabbitmqConfig struct {
	Address  string `json:"address" yaml:"address"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

// Config all config
type Config struct {
	DBCfg *DBConfig       `json:"DBService" yaml:"DBService"`
	ESCfg *[]ESConfig     `json:"EsService" yaml:"EsService"`
	GWCfg *Gateway        `json:"GatewayService" yaml:"GatewayService"`
	MQCfg *RabbitmqConfig `json:"RabbiqMQ" yaml:"RabbiqMQ"`
}

func TestJson(t *testing.T) {
	var cfg Config
	type test []struct {
		s1 interface{}
		s2 interface{}
	}
	err := ReadConfigFromJSONFile("config.json", &cfg)
	if err != nil {
		t.Error(err)
	}
	tests := test{
		{s1: cfg.DBCfg.Username, s2: "admin"},
		{s1: (*cfg.ESCfg)[0].Username, s2: "admin"},
		{s1: cfg.GWCfg.Count, s2: 3},
		{s1: cfg.MQCfg, s2: &RabbitmqConfig{
			Address:  "10.10.10.10",
			Port:     "5672",
			Username: "guest",
			Password: "guest",
		}},
	}

	for _, tc := range tests {
		if !reflect.DeepEqual(tc.s1, tc.s2) {
			t.Errorf("compare error, s1: %#v, s2: %#v", tc.s1, tc.s2)
		}
	}
}

func TestYAML(t *testing.T) {
	var cfg Config
	type test []struct {
		s1 interface{}
		s2 interface{}
	}
	err := ReadConfigFromYAMLFile("config.yaml", &cfg)
	if err != nil {
		t.Error(err)
	}
	tests := test{
		{s1: cfg.DBCfg.Username, s2: "admin"},
		{s1: (*cfg.ESCfg)[0].Username, s2: "admin"},
		{s1: cfg.GWCfg.Count, s2: 3},
		{s1: cfg.MQCfg, s2: &RabbitmqConfig{
			Address:  "10.10.10.10",
			Port:     "5672",
			Username: "guest",
			Password: "guest",
		}},
	}
	for _, tc := range tests {
		if !reflect.DeepEqual(tc.s1, tc.s2) {
			t.Errorf("compare error, s1: %#v, s2: %#v", tc.s1, tc.s2)
		}
	}
}
