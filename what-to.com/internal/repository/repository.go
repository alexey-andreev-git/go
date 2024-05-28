package repository

import (
	"what-to.com/internal/config"
	"what-to.com/internal/entity"
)

type (
	Repository interface {
		SetRepoConfig(config.ConfigT)
		GetRepoConfig() DBConfig
		GetRepoConfigStr() string
		CreateEntity(map[string]string) (entity.Entity, error)
		GetEntity(map[string]string) (entity.Entity, error)
		UpdateEntity(int, map[string]string) error
		DeleteEntity(int) error
		// GetEntityBy(interface{}) ([]entity.Entity, error)
		// GetEntityById(int) ([]entity.Entity, error)
		// GetEntityByName(string) ([]entity.Entity, error)
	}
)
