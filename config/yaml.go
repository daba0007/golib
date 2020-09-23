package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ReadConfigFromYAMLFile read config from json files
func ReadConfigFromYAMLFile(path string, conf interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return ReadConfigFromYAMLBytes(data, conf)
}

// ReadConfigFromYAMLBytes read config from bytes
func ReadConfigFromYAMLBytes(d []byte, conf interface{}) error {
	if err := yaml.Unmarshal(d, conf); err != nil {
		return err
	}
	return nil
}
