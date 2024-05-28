package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"what-to.com/internal/config"
	"what-to.com/internal/entity"
	"what-to.com/internal/resources"

	_ "github.com/lib/pq"
)

const (
	ENTITY_CREATE = `INSERT INTO entities (entity_name, entity_comment) VALUES ($1, $2) RETURNING entity_id`
	ENTITY_GET    = `SELECT entity_id, entity_name, entity_comment FROM entities`
	ENTITY_UPDATE = `UPDATE entities SET entity_name=$1, entity_comment=$2 WHERE entity_id=$3`
	ENTITY_DELETE = `DELETE FROM entities WHERE entity_id=$1`
)

type (
	// DBConfig holds the database connection configuration
	DBConfig struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}

	// PgRepository is the PostgreSQL repository
	PgRepository struct {
		DB        *sql.DB
		appConfig *config.Config
		dbConfig  DBConfig
	}
)

var entityFieldToColumn = map[string]string{
	"id":      "entity_id",
	"name":    "entity_name",
	"comment": "entity_comment",
}

// NewPgRepository initializes a new PostgreSQL repository
func NewPgRepository(appConfig *config.Config) *PgRepository {
	r := &PgRepository{
		DB:        nil,
		appConfig: appConfig,
	}
	r.SetRepoConfig(appConfig.GetConfig()["database"].(config.ConfigT))
	r.connectToDb()
	return r
}

// connectToDb connects to the PostgreSQL database
func (r *PgRepository) connectToDb() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password, r.dbConfig.DBName)

	r.checkDB()

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		r.appConfig.GetLogger().Fatal("Connection to the database failed:", err)
	}

	err = db.Ping()
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to execute ping on the database:", err)
	}

	r.appConfig.GetLogger().Info("PostgreSQL DB successfully connected!")
	r.DB = db
	r.UpdateDB()
}

// SetRepoConfig sets the DBConfig struct
func (r *PgRepository) SetRepoConfig(dbConfigP config.ConfigT) {
	r.dbConfig = DBConfig{
		Host:     dbConfigP["host"].(string),
		Port:     int(dbConfigP["port"].(int)),
		User:     dbConfigP["user"].(string),
		Password: dbConfigP["password"].(string),
		DBName:   dbConfigP["dbname"].(string),
	}
}

// GetRepoConfig returns the DBConfig
func (r *PgRepository) GetRepoConfig() DBConfig {
	return r.dbConfig
}

// GetRepoConfigStr returns the DBConfig as a string
func (r *PgRepository) GetRepoConfigStr() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password)
}

// checkDB checks if the database exists, and creates it if it doesn't
func (r *PgRepository) checkDB() bool {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	var exists int
	db.QueryRow("SELECT 1 FROM pg_database WHERE datname=$1", r.dbConfig.DBName).Scan(&exists)

	if exists == 0 {
		r.appConfig.GetLogger().Warn("Database does not exist. Creating...")
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", r.dbConfig.DBName))
		if err != nil {
			r.appConfig.GetLogger().Fatal("Failed to create database:", err)
		}
		r.appConfig.GetLogger().Info("Database created successfully.")
	}

	return true
}

// UpdateDB updates the database schema
func (r *PgRepository) UpdateDB() {
	appRes := resources.NewAppSources()
	fn := r.appConfig.GetConfig()[config.KeyInitDbFileName].(string)
	data, err := appRes.GetRes().ReadFile(fn)
	if err != nil {
		r.appConfig.GetLogger().Fatal("File read error [%s] "+fn, err)
	}
	_, err = r.DB.Exec(string(data))
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to update database:", err)
	}
}

// CreateEntity creates a new entity in the database
func (r *PgRepository) CreateEntity(ent map[string]string) (entity.Entity, error) {
	createdEntity := entity.Entity{
		Id:      0,
		Name:    ent["Name"],
		Comment: ent["Comment"],
	}
	err := r.DB.QueryRow(ENTITY_CREATE, ent["Name"], ent["Comment"]).Scan(&createdEntity.Id)
	return createdEntity, err
}

// GetEntity retrieves an entity from the database based on the given filter
func (r *PgRepository) GetEntity(filter map[string]string) (entity.Entity, error) {
	var conditions []string
	var args []interface{}
	i := 1
	for key, value := range filter {
		column, ok := entityFieldToColumn[strings.ToLower(key)]
		if !ok {
			return entity.Entity{}, fmt.Errorf("invalid field: %s", key)
		}
		conditions = append(conditions, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}
	query := ENTITY_GET
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return entity.Entity{}, err
	}
	defer rows.Close()

	var entities []entity.Entity
	for rows.Next() {
		var e entity.Entity
		if err := rows.Scan(&e.Id, &e.Name, &e.Comment); err != nil {
			return entity.Entity{}, err
		}
		entities = append(entities, e)
	}
	if len(entities) > 1 {
		return entity.Entity{}, fmt.Errorf("more than one entity found")
	}
	if len(entities) == 0 {
		return entity.Entity{}, fmt.Errorf("entity not found")
	}
	if err = rows.Err(); err != nil {
		return entity.Entity{}, err
	}
	return entities[0], nil
}

func (r *PgRepository) UpdateEntity(id int, ent map[string]string) error {
	_, err := r.DB.Exec(ENTITY_UPDATE, ent["Name"], ent["Comment"], id)
	return err
}

func (r *PgRepository) DeleteEntity(id int) error {
	_, err := r.DB.Exec(ENTITY_DELETE, id)
	return err
}
