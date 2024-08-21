package repository

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"what-to.com/internal/config"
	"what-to.com/internal/models"
)

type (
	GormPgRepository struct {
		db        *gorm.DB
		appConfig *config.Config
		dbConfig  DBConfig
	}
)

func NewGormPgRepository(conf *config.Config) *GormPgRepository {
	r := &GormPgRepository{
		appConfig: conf,
	}
	r.SetRepoConfig(r.appConfig.GetConfig()["database"].(config.ConfigT))
	r.connectToDb()
	return r
}

// connectToDb connects to the PostgreSQL database
func (r *GormPgRepository) connectToDb() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password, r.dbConfig.DBName)

	r.checkDB()

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		r.appConfig.GetLogger().Fatal("Connection to the database failed:", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to connect to database:", err)
	}

	err = sqlDb.Ping()
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to execute ping on the database:", err)
	}

	r.appConfig.GetLogger().Info("PostgreSQL DB successfully connected!")
	r.db = db
	r.UpdateDB()
}

// SetRepoConfig sets the DBConfig struct
func (r *GormPgRepository) SetRepoConfig(dbConfigP config.ConfigT) {
	r.dbConfig = DBConfig{
		Host:     dbConfigP["host"].(string),
		Port:     int(dbConfigP["port"].(int)),
		User:     dbConfigP["user"].(string),
		Password: dbConfigP["password"].(string),
		DBName:   dbConfigP["dbname"].(string),
	}
}

// GetRepoConfig returns the DBConfig
func (r *GormPgRepository) GetRepoConfig() DBConfig {
	return r.dbConfig
}

// GetRepoConfigStr returns the DBConfig as a string
func (r *GormPgRepository) GetRepoConfigStr() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password)
}

// checkDB checks if the database exists, and creates it if it doesn't
func (r *GormPgRepository) checkDB() bool {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to connect to database:", err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to connect to database:", err)
	}
	defer sqlDb.Close()

	var exists int
	// sqlDb.ExecQueryRow("SELECT 1 FROM pg_database WHERE datname=$1", r.dbConfig.DBName).Scan(&exists)
	db.Raw("SELECT 1 FROM pg_database WHERE datname=?", r.dbConfig.DBName).Scan(&exists)

	if exists == 0 {
		r.appConfig.GetLogger().Warn("Database does not exist. Creating...")
		_, err := sqlDb.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", r.dbConfig.DBName))
		if err != nil {
			r.appConfig.GetLogger().Fatal("Failed to create database:", err)
		}
		r.appConfig.GetLogger().Info("Database created successfully.")
	}

	return true
}

// UpdateDB updates the database schema
func (r *GormPgRepository) UpdateDB() {
	// r.db.AutoMigrate(
	// 	&models.User{},
	// 	&models.Person{},
	// 	&models.Company{},
	// 	&models.Address{},
	// 	&models.Email{},
	// 	&models.Phone{},
	// 	&models.CompanyAddress{},
	// 	&models.CompanyEmail{},
	// 	&models.CompanyPhone{},
	// 	&models.PersonAddress{},
	// 	&models.PersonEmail{},
	// 	&models.PersonPhone{},
	// )
	r.db.AutoMigrate(
		&models.Address{},
	)
}

func (r *GormPgRepository) Create(ctx context.Context, value interface{}) error {
	return r.db.WithContext(ctx).Create(value).Error
}

func (r *GormPgRepository) Update(ctx context.Context, value interface{}) error {
	return r.db.WithContext(ctx).Save(value).Error
}

func (r *GormPgRepository) Delete(ctx context.Context, value interface{}) error {
	return r.db.WithContext(ctx).Delete(value).Error
}

func (r *GormPgRepository) FindByID(ctx context.Context, id uint, out interface{}) error {
	return r.db.WithContext(ctx).First(out, id).Error
}

func (r *GormPgRepository) FindAll(ctx context.Context, out interface{}) error {
	return r.db.WithContext(ctx).Find(out).Error
}
