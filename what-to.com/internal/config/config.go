package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v2"
	"what-to.com/internal/logger"
)

type ConfigT = map[interface{}]interface{}

var (
	envConfigFile string = "pg_db_connection.yaml"
)

type Config struct {
	customLogger logger.Logger
	configFile   string
	Config       ConfigT
}

func init() {
	// Define and parse command-line arguments
	flag.StringVar(&envConfigFile, "config", envConfigFile, "path to the configuration file")
	flag.Parse()

	// Override configFile if an environment variable is set
	if configFile := os.Getenv("WHATTO_CONFIG_FILE_PATH"); configFile != "" {
		envConfigFile = configFile
	}
}

func NewConfig() *Config {
	c := &Config{
		configFile:   envConfigFile,
		Config:       nil,
		customLogger: logger.NewCustomLogger("whattoapp.log"),
	}
	// c.log = logger.NewCustomLogger(c.configFile)
	c.ReadConfig()
	return c
}

func (c *Config) GetConfig() ConfigT {
	return c.Config
}
func (c *Config) GetLogger() logger.Logger {
	return c.customLogger
}

func (c *Config) ReadConfig() {
	// Read the YAML file
	yamlFile, err := os.ReadFile(c.configFile)
	if err != nil {
		c.customLogger.Fatal("Error reading the configuration file:", err)
	}

	// Parse the YAML file into a map
	err = yaml.Unmarshal(yamlFile, &c.Config)
	if err != nil {
		c.customLogger.Fatal("Parsing YAML file error:", err)
	}

	c.customLogger.Info("Config successfully read!")
}
