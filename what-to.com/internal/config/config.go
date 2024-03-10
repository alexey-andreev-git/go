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
	Config     ConfigT
)

func ReadConfig() {
	// Read the YAML file
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error reading the configuration file:", err)
	}

	// Parse the YAML file into a map
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatal("Parsing YAML file error:", err)
	}

	log.Info("Config successfully read!")
}
