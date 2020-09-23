package config

import (
	"encoding/json"
	"io/ioutil"
)

// ReadConfigFromJSONFile read config from json files
func ReadConfigFromJSONFile(path string, conf interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return ReadConfigFromJSONBytes(data, conf)
}

// ReadConfigFromJSONBytes read config from bytes
func ReadConfigFromJSONBytes(d []byte, conf interface{}) error {
	if err := json.Unmarshal(d, conf); err != nil {
		return err
	}
	return nil
}
