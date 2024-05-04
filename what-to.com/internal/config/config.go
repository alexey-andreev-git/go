package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v2"
	"what-to.com/internal/logger"
)

type ConfigT = map[interface{}]interface{}

const (
	InitDbFileName        = "appfs/sql/initdb.sql"
	KeyInitDbFileName     = "initDbFileName" //Key in the config map
	InitConfigFileName    = "pg_db_connection.yaml"
	KeyInitConfigFileName = "configFileName" //Key in the config map
	InitLogFileName       = "whattoapp.log"
	KeyInitLogFileName    = "logFileName" //Key in the config map
)

var (
	envConfigFile string = InitConfigFileName // Default value from const
	envLogFile    string = InitLogFileName    // Default value from const
)

type Config struct {
	customLogger logger.Logger
	configFile   string
	cConfig      ConfigT
}

func init() {
	// Define and parse command-line arguments
	flag.StringVar(&envConfigFile, "config", envConfigFile, "path to the configuration file")
	flag.Parse()

	flag.StringVar(&envLogFile, "log", envLogFile, "path to the log file")
	flag.Parse()

	// Override configFile if an environment variable is set
	if envConfigFile == "" {
		// Checking the environment variable in case the flag is not set
		if configFile := os.Getenv("WHATTO_CONFIG_FILE_PATH"); configFile != "" {
			envConfigFile = configFile
		}
	}
	if envLogFile == "" {
		// Checking the environment variable in case the flag is not set
		if logFile := os.Getenv("WHATTO_LOG_FILE_PATH"); logFile != "" {
			envLogFile = logFile
		}
	}
}

func NewConfig() *Config {
	c := &Config{
		configFile:   envConfigFile,
		cConfig:      nil,
		customLogger: logger.NewCustomLogger(envLogFile),
	}
	c.ReadConfig()
	return c
}

func (c *Config) GetConfig() ConfigT {
	return c.cConfig
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
	err = yaml.Unmarshal(yamlFile, &c.cConfig)
	if err != nil {
		c.customLogger.Fatal("Parsing YAML file error:", err)
	}

	// Add calculated parameters to the YAML file into a map
	c.cConfig[KeyInitConfigFileName] = c.configFile
	c.cConfig[KeyInitDbFileName] = InitDbFileName
	c.cConfig[KeyInitLogFileName] = envLogFile

	// Report the YAML map created
	c.customLogger.Info("Config successfully read!")
}
