package repository

import (
	"what-to.com/internal/config"
	"what-to.com/internal/entity"
)

type (
	Repository interface {
		SetRepoConfig(config.ConfigT)
		GetRepoConfig() config.ConfigT
		GetRepoConfigStr() string
		CreateEntity(*entity.Entity) error
	}
)
