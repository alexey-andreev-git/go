package repository

import "what-to.com/internal/entity"

type Repository interface {
	CreateEntity(*entity.Entity) error
}
