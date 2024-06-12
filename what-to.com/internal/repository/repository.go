package repository

import (
	"database/sql"

	"what-to.com/internal/config"
	"what-to.com/internal/entity"
)

type (
	Repository interface {
		SetRepoConfig(config.ConfigT)
		GetRepoConfig() DBConfig
		GetRepoConfigStr() string

		CreateEntity(map[string]interface{}) (entity.Entity, error)
		GetEntity(map[string]interface{}) (entity.Entity, error)
		UpdateEntity(map[string]interface{}) (sql.Result, error)
		DeleteEntity(map[string]interface{}) (sql.Result, error)

		CreateEntityData(map[string]interface{}) ([]entity.EntityData, error)
		GetEntityData(map[string]interface{}) ([]entity.EntityData, error)
		UpdateEntityData(map[string]interface{}) (sql.Result, error)
		DeleteEntityData(map[string]interface{}) (sql.Result, error)

		CreateEntityDataRef(map[string]interface{}) ([]entity.EntityDataReference, error)
		GetEntityDataRef(map[string]interface{}) ([]entity.EntityDataReference, error)
		UpdateEntityDataRef(map[string]interface{}) (sql.Result, error)
		DeleteEntityDataRef(map[string]interface{}) (sql.Result, error)
	}
)
