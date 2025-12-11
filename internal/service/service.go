package service

import "cargo-m/internal/repository"

type BaseService[T any] struct {
	baseRepo repository.IRepository[T]
}

func NewBaseService[E any](repo repository.IRepository[E]) *BaseService[E] {
	return &BaseService[E]{baseRepo: repo}
}
