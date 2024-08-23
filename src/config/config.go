package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Port          string `yaml:"port"`
	Storage       string `yaml:"storage"`
	FilerPath     string `yaml:"path"`
	WebResources  string `yaml:"resources"`
	FfmpegPath    string `yaml:"ffmpeg"`
	UrlPublicKeys string `yaml:"url_public_keys"`
}

func (c Config) check() bool {
	return true
}

func ReadConfig(path string) (*Config, error) {
	if f, err := os.Open(path); err == nil {
		decoder := yaml.NewDecoder(f)
		conf := &Config{Port: "9006"}
		if err := decoder.Decode(conf); err == nil {
			if conf.check() {
				return conf, nil
			}
			return nil, errors.New("bad configuration")
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
