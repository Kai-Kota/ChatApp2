package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Kai-Kota/ChatApp2/backend/types"
)

var (
	ErrJsonUnmarshal     = errors.New("failed to parse product from request body")
	ErrProductIdMismatch = errors.New("product ID in path does not match product ID in body")
)

type Products struct {
	store types.Store
}

func NewProductsDomain(s types.Store) *Products {
	return &Products{
		store: s,
	}
}

func (d *Products) PutProduct(ctx context.Context, id string, body []byte) (*types.Product, error) {
	product := types.Product{}
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, fmt.Errorf("%v", ErrJsonUnmarshal)
	}

	if product.Id != id {
		return nil, fmt.Errorf("%v", ErrProductIdMismatch)
	}

	err := d.store.Put(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &product, nil
}

func (d *Products) GetProduct(ctx context.Context, id string) (*types.Product, error) {
	product, err := d.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return product, nil
}
