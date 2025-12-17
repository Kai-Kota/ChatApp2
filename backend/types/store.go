package types

import "context"

type Store interface {
	Put(context.Context, Product) error
}
