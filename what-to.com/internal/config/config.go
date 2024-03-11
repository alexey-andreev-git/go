package config

import (
	"os"

	"gopkg.in/yaml.v2"
	"what-to.com/internal/logger"
)

type ConfigT = map[interface{}]interface{}

var (
	log        = logger.CustomLogger
	configFile = "pg_db_connection.yaml"
	// Config     ConfigT
)

type Config struct {
	configFile string // = "pg_db_connection.yaml"
	Config     ConfigT
}

func NewConfig(confFileName string) *Config {
	c := &Config{
		configFile: confFileName,
		Config:     nil,
	}
	c.ReadConfig()
	return c
}

func (c *Config) GetConfig() *ConfigT {
	return &c.Config
}

func (c *Config) ReadConfig() {
	// Read the YAML file
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error reading the configuration file:", err)
	}

	// Parse the YAML file into a map
	err = yaml.Unmarshal(yamlFile, &c.Config)
	if err != nil {
		log.Fatal("Parsing YAML file error:", err)
	}

	log.Info("Config successfully read!")
}
