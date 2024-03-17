package repository

import (
	"database/sql"
	"fmt"

	"what-to.com/internal/config"
	"what-to.com/internal/entity"

	_ "github.com/lib/pq"
)

type (
	// create a structure to hold the database connection
	DBConfig struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}

	PgRepository struct {
		DB        *sql.DB
		appConfig *config.Config
		dbConfig  DBConfig
	}
)

func NewPgRepository(appConfig *config.Config) *PgRepository {
	r := &PgRepository{
		DB:        nil,
		appConfig: appConfig,
	}
	r.SetRepoConfig(appConfig.GetConfig()["database"].(config.ConfigT))
	r.connectToDb()
	return r
}

// ConnectToDB connects to the database
func (r *PgRepository) connectToDb() {
	// Сначала подключаемся к базе данных `postgres` для проверки существования целевой БД
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to connect to database:", err)
	}

	// Проверяем, существует ли целевая база данных
	var exists int
	db.QueryRow("SELECT 1 FROM pg_database WHERE datname=$1", r.dbConfig.DBName).Scan(&exists)

	if exists == 0 {
		// База данных не существует, создаем ее
		r.appConfig.GetLogger().Warn("Database does not exists. Creating...")
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", r.dbConfig.DBName))
		if err != nil {
			r.appConfig.GetLogger().Fatal("Failed to create database:", err)
		}
		r.appConfig.GetLogger().Info("Database created successfully.")
	}

	// Закрываем соединение с базой данных `postgres`
	db.Close()

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password, r.dbConfig.DBName)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		r.appConfig.GetLogger().Fatal("Connection to the database failed:", err)
	}

	err = db.Ping()
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to execute ping on the database:", err)
	}

	r.appConfig.GetLogger().Info("PostgreSQL DB successfully connected!")
	r.DB = db
}

// SetDBConfig sets the DBConfig struct
func (r *PgRepository) SetRepoConfig(dbConfigP config.ConfigT) {
	r.dbConfig = DBConfig{
		Host:     dbConfigP["host"].(string),
		Port:     int(dbConfigP["port"].(int)),
		User:     dbConfigP["user"].(string),
		Password: dbConfigP["password"].(string),
		DBName:   dbConfigP["dbname"].(string),
	}
}

func (r *PgRepository) GetRepoConfigStr() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password)
}

func (r *PgRepository) GetRepoConfig() DBConfig {
	return r.dbConfig
}

// Пример функции для добавления новой сущности в базу данных
func (r *PgRepository) CreateEntity(ye *entity.Entity) error {
	query := `INSERT INTO company_entity_table (name) VALUES ($1)`
	_, err := r.DB.Exec(query, ye.Name)
	return err
}
