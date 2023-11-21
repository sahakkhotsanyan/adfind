package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WordLists map[string]string `yaml:"word_lists"` // map[websiteType]wordList.txt
}

func NewConfig(location string) (*Config, error) {
	cfg := Config{}
	file, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&cfg)

	return &cfg, err
}

func (c *Config) GetWordListFileName(websiteType string) string {
	return c.WordLists[websiteType]
}
