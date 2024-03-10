package repository

import (
	"database/sql"
	"fmt"

	"what-to.com/internal/config"
	"what-to.com/internal/logger"

	_ "github.com/lib/pq"
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

var (
	dbConfig = DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "your_username",
		Password: "your_password",
		DBName:   "your_dbname",
	}
	log = logger.CustomLogger
)

func ConnectToDB() *sql.DB {
	// Сначала подключаемся к базе данных `postgres` для проверки существования целевой БД
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Проверяем, существует ли целевая база данных
	var exists int
	db.QueryRow("SELECT 1 FROM pg_database WHERE datname=$1", dbConfig.DBName).Scan(&exists)

	if exists == 0 {
		// База данных не существует, создаем ее
		log.Warn("Database does not exists. Creating...")
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbConfig.DBName))
		if err != nil {
			log.Fatal("Failed to create database:", err)
		}
		log.Info("Database created successfully.")
	}

	// Закрываем соединение с базой данных `postgres`
	db.Close()

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Connection to the database failed:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to execute ping on the database:", err)
	}

	log.Info("PostgreSQL DB successfully connected!")
	return db
}

// SetDBConfig sets the DBConfig struct
func SetDBConfig(dbConfigP config.ConfigT) {
	dbConfig = DBConfig{
		Host:     dbConfigP["host"].(string),
		Port:     int(dbConfigP["port"].(int)),
		User:     dbConfigP["user"].(string),
		Password: dbConfigP["password"].(string),
		DBName:   dbConfigP["dbname"].(string),
	}

}

// // Read pg_db_connection.yaml and return a DBConfig struct
// func ReadDBConfig() DBConfig {
// 	// Read the file and return a DBConfig struct
// 	// Assuming the YAML file has the following structure:
// 	// host: localhost
// 	// port: 5432
// 	// user: your_username
// 	// password: your_password
// 	// dbname: your_dbname

// 	// Read the YAML file
// 	yamlFile, err := os.ReadFile("pg_db_connection.yaml")
// 	if err != nil {
// 		log.Fatal("Error reading the configuration file:", err)
// 	}

// 	// Parse the YAML file into a map
// 	var yamlData map[string]interface{}
// 	err = yaml.Unmarshal(yamlFile, &yamlData)
// 	if err != nil {
// 		log.Fatal("Parsing YAML file error:", err)
// 	}

// 	// Fill the DBConfig struct with the values from the YAML file
// 	var dbConnectConfig = yamlData["database"].(map[interface{}]interface{})
// 	config := DBConfig{
// 		Host:     dbConnectConfig["host"].(string),
// 		Port:     int(dbConnectConfig["port"].(int)),
// 		User:     dbConnectConfig["user"].(string),
// 		Password: dbConnectConfig["password"].(string),
// 		DBName:   dbConnectConfig["dbname"].(string),
// 	}

// 	return config
// }
