package repository

import (
	"database/sql"
	"fmt"

	"what-to.com/internal/config"
	"what-to.com/internal/entity"
	"what-to.com/internal/logger"

	_ "github.com/lib/pq"
)

type PgRepository struct {
	DB *sql.DB
}

func NewPgRepository(appConfig *map[interface{}]interface{}) *PgRepository {
	r := &PgRepository{DB: nil}
	r.SetRepoConfig((*appConfig)["database"].(config.ConfigT))
	r.ConnectToRepo()
	return r
}

// Пример функции для добавления новой сущности в базу данных
func (r *PgRepository) CreateEntity(ye *entity.Entity) error {
	query := `INSERT INTO company_entity_table (name) VALUES ($1)`
	_, err := r.DB.Exec(query, ye.Name)
	return err
}

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

func (r *PgRepository) ConnectToRepo() {
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
	r.DB = db
}

// SetDBConfig sets the DBConfig struct
func (r *PgRepository) SetRepoConfig(dbConfigP config.ConfigT) {
	dbConfig = DBConfig{
		Host:     dbConfigP["host"].(string),
		Port:     int(dbConfigP["port"].(int)),
		User:     dbConfigP["user"].(string),
		Password: dbConfigP["password"].(string),
		DBName:   dbConfigP["dbname"].(string),
	}
}

func (r *PgRepository) GetRepoConfigStr() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password)
}

func (r *PgRepository) GetRepoConfig() DBConfig {
	return dbConfig
}
