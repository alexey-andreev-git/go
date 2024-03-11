package repository

import (
	"what-to.com/internal/config"
	"what-to.com/internal/entity"
)

type (
	Repository interface {
		SetRepoConfig(config.ConfigT)
		GetRepoConfigStr() string
		CreateEntity(*entity.Entity) error
	}
	DataParamsT = map[interface{}]interface{}
)
