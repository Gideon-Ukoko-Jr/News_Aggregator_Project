package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	UserServicePreferenceURL    string `yaml:"user_service_fetch_preference_url"`
	NewsAggregatorRecentNewsURL string `yaml:"news_aggregator_recent_news_url"`
	SpecialKey                  string `yaml:"special_key"`
}

func LoadConfig() (*Config, error) {
	// Getting the absolute path to the config.yaml file
	configPath, err := filepath.Abs("config/config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
