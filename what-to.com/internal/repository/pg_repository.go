package repository

import (
	"database/sql"

	"what-to.com/internal/entity"
)

type PgRepository struct {
	DB *sql.DB
}

func NewPgRepository(db *sql.DB) *PgRepository {
	return &PgRepository{DB: db}
}

// Пример функции для добавления новой сущности в базу данных
func (r *PgRepository) CreateEntity(ye *entity.Entity) error {
	query := `INSERT INTO company_entity_table (name) VALUES ($1)`
	_, err := r.DB.Exec(query, ye.Name)
	return err
}
