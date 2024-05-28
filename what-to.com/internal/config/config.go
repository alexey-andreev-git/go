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
	KeyInitDbFileName     = "initDbFileName"
	InitConfigFileName    = "pg_db_connection.yaml"
	KeyInitConfigFileName = "configFileName"
	InitLogFileName       = "whattoapp.log"
	KeyInitLogFileName    = "logFileName"
)

var (
	envConfigFile string = InitConfigFileName
	envLogFile    string = InitLogFileName
)

type Config struct {
	customLogger logger.Logger
	configFile   string
	cConfig      ConfigT
}

func init() {
	flag.StringVar(&envConfigFile, "config", envConfigFile, "path to the configuration file")
	flag.StringVar(&envLogFile, "log", envLogFile, "path to the log file")
	flag.Parse()

	if envConfigFile == "" {
		if configFile := os.Getenv("WHATTO_CONFIG_FILE_PATH"); configFile != "" {
			envConfigFile = configFile
		}
	}
	if envLogFile == "" {
		if logFile := os.Getenv("WHATTO_LOG_FILE_PATH"); logFile != "" {
			envLogFile = logFile
		}
	}
}

func NewConfig() *Config {
	c := &Config{
		configFile:   envConfigFile,
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
	yamlFile, err := os.ReadFile(c.configFile)
	if err != nil {
		c.customLogger.Fatal("Error reading the configuration file:", err)
	}

	err = yaml.Unmarshal(yamlFile, &c.cConfig)
	if err != nil {
		c.customLogger.Fatal("Parsing YAML file error:", err)
	}

	c.cConfig[KeyInitConfigFileName] = c.configFile
	c.cConfig[KeyInitDbFileName] = InitDbFileName
	c.cConfig[KeyInitLogFileName] = envLogFile

	c.customLogger.Info("Config successfully read!")
}
