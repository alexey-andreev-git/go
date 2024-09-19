package repository

import (
	// "database/sql"

	"context"

	"what-to.com/internal/config"
	// "what-to.com/internal/entity"
)

type (
	Repository interface {
		SetRepoConfig(config.ConfigT)
		GetRepoConfig() DBConfig
		GetRepoConfigStr() string

		Create(ctx context.Context, value interface{}) error
		Update(ctx context.Context, value interface{}) error
		Delete(ctx context.Context, value interface{}) error
		FindByID(ctx context.Context, entId uint, id uint) (interface{}, error)
		FindAll(ctx context.Context, entId uint) (interface{}, error)

		// CreateEntity(map[string]interface{}) (entity.Entity, error)
		// GetEntity(map[string]interface{}) (entity.Entity, error)
		// UpdateEntity(map[string]interface{}) (sql.Result, error)
		// DeleteEntity(map[string]interface{}) (sql.Result, error)

		// CreateEntityData(map[string]interface{}) ([]entity.EntityData, error)
		// GetEntityData(map[string]interface{}) ([]entity.EntityData, error)
		// UpdateEntityData(map[string]interface{}) (sql.Result, error)
		// DeleteEntityData(map[string]interface{}) (sql.Result, error)

		// CreateEntityDataRef(map[string]interface{}) ([]entity.EntityDataReference, error)
		// GetEntityDataRef(map[string]interface{}) ([]entity.EntityDataReference, error)
		// UpdateEntityDataRef(map[string]interface{}) (sql.Result, error)
		// DeleteEntityDataRef(map[string]interface{}) (sql.Result, error)
	}
)
