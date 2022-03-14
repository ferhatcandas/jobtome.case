package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func LoadYMLConfig(file string, cfg interface{}) error {
	f, err := os.Open(file)

	if err != nil {
		return err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}
