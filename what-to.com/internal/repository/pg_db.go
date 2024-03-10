package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"what-to.com/internal/logger"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

// ConnectToDB connects to the database
// create a structure to hold the database connection
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

var config = DBConfig{
	Host:     "localhost",
	Port:     5432,
	User:     "your_username",
	Password: "your_password",
	DBName:   "your_dbname",
}

func ConnectToDB() *sql.DB {
	// Сначала подключаемся к базе данных `postgres` для проверки существования целевой БД
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Проверяем, существует ли целевая база данных
	var exists int
	db.QueryRow("SELECT 1 FROM pg_database WHERE datname=$1", config.DBName).Scan(&exists)

	if exists == 0 {
		// База данных не существует, создаем ее
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", config.DBName))
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		fmt.Println("Database created successfully")
	}

	// Закрываем соединение с базой данных `postgres`
	db.Close()

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Log.Fatalf("Connection to the database failed: %v", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Log.Fatalf("Failed to execute ping on the database: %v", err)
	}

	fmt.Println("PostgreSQL DB successfully connected!")
	return db
}

// SetDBConfig sets the DBConfig struct
func SetDBConfig(dbConfig DBConfig) {
	config = dbConfig
}

// Read pg_db_connection.yaml and return a DBConfig struct
func ReadDBConfig() DBConfig {
	// Read the file and return a DBConfig struct
	// Assuming the YAML file has the following structure:
	// host: localhost
	// port: 5432
	// user: your_username
	// password: your_password
	// dbname: your_dbname

	// Read the YAML file
	yamlFile, err := os.ReadFile("pg_db_connection.yaml")
	if err != nil {
		logger.Log.Fatalf("Error reading the configuration file: %v", err)
	}

	// Parse the YAML file into a map
	var yamlData map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &yamlData)
	if err != nil {
		logger.Log.Fatalf("Parsing YAML file error: %v", err)
	}

	// Fill the DBConfig struct with the values from the YAML file
	var dbConnectConfig = yamlData["database"].(map[interface{}]interface{})
	config := DBConfig{
		Host:     dbConnectConfig["host"].(string),
		Port:     int(dbConnectConfig["port"].(int)),
		User:     dbConnectConfig["user"].(string),
		Password: dbConnectConfig["password"].(string),
		DBName:   dbConnectConfig["dbname"].(string),
	}

	return config
}
