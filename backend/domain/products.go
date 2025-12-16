package domain

import (
	"ChatApp2/types"
)

type Products struct {
	store types.Store
}

func NewProductsDomain(s types.Store) *Products {
	return &Products{
		store: s,
	}
}
