package repository

import (
	"database/sql"

	"what-to/internal/entity"
)

type PgRepository struct {
	DB *sql.DB
}

func NewYourRepository(db *sql.DB) *YourRepository {
	return &YourRepository{DB: db}
}

// Пример функции для добавления новой сущности в базу данных
func (r *YourRepository) CreateYourEntity(ye entity.YourEntity) error {
	query := `INSERT INTO your_entity_table (name) VALUES ($1)`
	_, err := r.DB.Exec(query, ye.Name)
	return err
}

// Добавьте здесь другие функции репозитория
